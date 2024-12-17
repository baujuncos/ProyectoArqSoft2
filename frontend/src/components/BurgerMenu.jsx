import { useState, useEffect } from "react";
import {
    Box,
    IconButton,
    Drawer,
    DrawerOverlay,
    DrawerContent,
    DrawerCloseButton,
    DrawerHeader,
    DrawerBody,
    VStack,
    Button,
    useDisclosure
} from "@chakra-ui/react";
import { HamburgerIcon } from "@chakra-ui/icons";
import Cookies from "js-cookie";
import Popup from "./PopUp.jsx";
import PopupCreate from "./PopUpCreate.jsx";
import PopupRegister from "./PopUpRegister.jsx";
import ItemList from './ItemList';
import PopUpMicroservicios from "./PopUpMicroservicios.jsx";

const BurgerMenu = ({ onLogout }) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const { isOpen: isPopupOpen, onOpen: onOpenPopup, onClose: onClosePopup } = useDisclosure();
    const { isOpen: isPopupOpenCreate, onOpen: onOpenPopupCreate, onClose: onClosePopupCreate } = useDisclosure();
    const { isOpen: isPopupOpenRegister, onOpen: onOpenPopupRegister, onClose: onClosePopupRegister } = useDisclosure();
    const { isOpen: isPopupOpenMicroservicios, onOpen: onOpenPopupMicroservicios, onClose: onClosePopupMicroservicios } = useDisclosure();
    const [userId, setUserId] = useState(null);
    const [admin, setAdmin] = useState(null);
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        const storedUserId = Cookies.get('user_id');
        const storedAdmin = Cookies.get('admin');

        if (storedUserId) {
            setUserId(parseInt(storedUserId, 10));
        }

        if (storedAdmin) {
            setAdmin(storedAdmin === "1");
        }
    }, [isPopupOpen]);

    const handleLogout = () => {
        Cookies.remove('user_id');
        Cookies.remove('email');
        Cookies.remove('token');
        Cookies.remove('admin');
        setUserId(null);
        onLogout();
    };

    const fetchCourseDetails = async (courseId) => {
        try {
            const response = await fetch(`http://localhost:8081/courses/${courseId}`);
            if (!response.ok) {
                if (response.status === 404) {
                    console.warn(`Curso con ID ${courseId} no encontrado.`);
                    return null;
                }
                throw new Error(`Error del servidor para el curso con ID ${courseId}: ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error(`Error al obtener detalles del curso con ID ${courseId}:`, error);
            return null;
        }
    };

    const handleMyCourses = async () => {
        setLoading(true);
        setError(null);

        if (!userId) {
            setError('El usuario no está autenticado.');
            setLoading(false);
            return;
        }

        const url = `http://localhost:8080/inscripciones/usuario/${userId}`;

        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`Error del servidor: ${response.statusText}`);
            }

            const inscripciones = await response.json();
            const courseDetailsPromises = inscripciones.map((id) =>
                fetchCourseDetails(id).catch((error) => {
                    console.error(`Error al obtener detalles del curso ${id}:`, error);
                    return null;
                })
            );

            const courseDetails = await Promise.all(courseDetailsPromises);
            const validCourses = courseDetails.filter((course) => course !== null);
            setCourses(validCourses);

        } catch (error) {
            console.error('Error al obtener los cursos del usuario:', error);
            setError('No se pudieron obtener los cursos. Intenta nuevamente más tarde.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Box p={4}>
            <IconButton
                aria-label="Open Menu"
                icon={<HamburgerIcon />}
                onClick={onOpen}
                size="lg"
            />
            <Drawer placement="right" onClose={onClose} isOpen={isOpen}>
                <DrawerOverlay />
                <DrawerContent>
                    <DrawerCloseButton />
                    <DrawerHeader style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Menú</DrawerHeader>
                    <DrawerBody>
                        <VStack spacing={4}>
                            {userId ? (
                                <>
                                    <Button w="100%" onClick={handleLogout} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Cerrar sesión</Button>
                                    {admin ? (
                                        <>
                                            <Button w="100%" onClick={onOpenPopupCreate} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Crear curso</Button>
                                            <Button w="100%" onClick={onOpenPopupMicroservicios} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Gestión de Microservicios</Button>
                                        </>
                                    ) : (
                                        <>
                                            <Button w="100%" onClick={handleMyCourses} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Mis cursos</Button>
                                            {loading && <p>Cargando cursos...</p>}
                                            {courses.length > 0 && (
                                                <Box w="100%">
                                                    <ItemList courses={courses} bandera={1} />
                                                </Box>
                                            )}
                                            {error && <Box color="red.500">{error}</Box>}
                                        </>
                                    )}
                                </>
                            ) : (
                                <>
                                    <Button w="100%" onClick={onOpenPopup} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Iniciar sesión</Button>
                                    <Button w="100%" onClick={onOpenPopupRegister} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>Registrarme</Button>
                                </>
                            )}
                        </VStack>
                    </DrawerBody>
                </DrawerContent>
            </Drawer>
            <Popup isOpen={isPopupOpen} onClose={onClosePopup} />
            <PopupCreate isOpen={isPopupOpenCreate} onClose={onClosePopupCreate} />
            <PopupRegister isOpen={isPopupOpenRegister} onClose={onClosePopupRegister} />
            <PopUpMicroservicios isOpen={isPopupOpenMicroservicios} onClose={onClosePopupMicroservicios} />
        </Box>
    );
};

export default BurgerMenu;
