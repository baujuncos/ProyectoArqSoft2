import React, { useEffect, useState } from "react";
import axios from "axios";
import Cookies from "js-cookie";
import { Button, Stack, Card, CardBody, CardFooter, Text, Image, useDisclosure } from "@chakra-ui/react";
import Inscribirmebutton from "./Inscribirmebutton.jsx";
import EliminarButton from "./EliminarButton.jsx";
import PopupEdit from "./PopUpEdit.jsx";
import '../estilos/Inscribirmebutton.css';
import '../estilos/Course.css';
import PopupValorar from "./PopUpValorar.jsx"
import PopupSubirArchivo from "./PopUpArchivo.jsx";
import PopupSeeReview from "./PopUpSeeReview.jsx";
import PopupDetail from "./PopUpDetail.jsx";

const Item = ({ course, bandera }) => {
    const [userId, setUserId] = useState(null);
    const [isAdmin, setIsAdmin] = useState(false);
    const [isEnrolled, setIsEnrolled] = useState(false);
    const { isOpen: isPopupOpenEdit, onOpen: onOpenPopupEdit, onClose: onClosePopupEdit } = useDisclosure();
    const { isOpen: isPopupOpenValorar, onOpen: onOpenPopupValorar, onClose: onClosePopupValorar } = useDisclosure();
    const { isOpen: isPopupOpenSubirArchivo, onOpen: onOpenPopupSubirArchivo, onClose: onClosePopupSubirArchivo } = useDisclosure();
    const { isOpen: isPopupOpenSeeReview , onOpen: onOpenPopupSeeReview , onClose: onClosePopupSeeReview } = useDisclosure();
    const { isOpen: isPopupOpenDetail, onOpen: onOpenPopupDetail, onClose: onClosePopupDetail } = useDisclosure();


    const formattedDate = new Date(course.fecha_inicio).toLocaleDateString('es-ES', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    });

    useEffect(() => {
        const storedUserId = Cookies.get('user_id');
        if (storedUserId) {
            setUserId(parseInt(storedUserId, 10));
        }

        const storedAdmin = Cookies.get('admin');
        if (storedAdmin) {
            setIsAdmin(storedAdmin === "1");
        }
    }, []);

    useEffect(() => {
        const checkEnrollment = async () => {
            if (userId && course.course_id) { // Asegúrate de que ambos valores estén definidos
                try {
                    // Realiza la solicitud al backend
                    const response = await axios.get(`http://localhost:8080/inscripciones/usuario/${userId}`);
                    const inscripciones = response.data; // Esto debería ser un arreglo de IDs (strings)

                    // Verifica si el ID del curso está en el arreglo de inscripciones
                    const enrolled = inscripciones.includes(course.course_id.toString()); // Convertimos a string para asegurarnos
                    setIsEnrolled(enrolled);
                } catch (error) {
                    console.error('Error checking enrollment:', error);
                }
            }
        };

        checkEnrollment();
    }, [userId, course.course_id]);


    const getProfesorName = (profesor_id) => {
        const profesores = {
            2: 'Juan Lopez',
            4: 'Margarita de Marcos',
            8: 'Gustavo Jacobo',
            17: 'Rodolfo Perez',
            18: 'Sebastian Colidio',
            19: 'Lucas Beltran'
        };
        return profesores[profesor_id] || 'Profesor desconocido';
    };

    const handleEditCourse = () => {
        onOpenPopupEdit();
    };

    const handleValorarCourse =()=>{
        onOpenPopupValorar();
    }

    const handleSubirArchivo = () => {
        onOpenPopupSubirArchivo();
    };

    const handleSeeReview = () => {
        onOpenPopupSeeReview();
    };


    return (
        <Card direction={{ base: 'column', sm: 'row' }} overflow='hidden' variant='outline'>
            {bandera !== 1 ? (
                <Image
                    objectFit='cover'
                    maxW={{ sm: '250px' }}
                    src={course.url_image}
                    alt='Imagen Curso'
                />
            ) : null}

            <Stack>
                <CardBody className='body'>
                    <h1 style={{ fontFamily: 'Spoof Trial, sans-serif', fontWeight: 800, fontSize: 30 }}>{course.nombre}</h1>

                    <Text py='2' className="card-text">{course.descripcion}</Text>
                    <Text marginBottom='3px' display='flex' py='2' alignItems='center' className="card-text">
                        <img src="/estrella.png" alt="estrella" width="20px" height="20px" style={{ marginRight: '5px' }} />
                        {course.valoracion}/5
                    </Text>
                    <Button onClick={onOpenPopupDetail}>Ver detalle</Button>
                </CardBody>
                <CardFooter>
                    {userId && (
                        isAdmin ? (
                            <>

                            </>
                        ) : (
                            isEnrolled ? (
                                bandera !== 1 ? (
                                    <>

                                    </>
                                ) : null
                            ) : (
                                bandera !== 1 ? (
                                    <>

                                    </>
                ) : null
                            )
                        )
                    )}

                </CardFooter>
            </Stack>
            <PopupEdit isOpen={isPopupOpenEdit} onClose={onClosePopupEdit} courseId={course.course_id} />
            <PopupValorar isOpen={isPopupOpenValorar} onClose={onClosePopupValorar} courseId={course.course_id} />
            <PopupSubirArchivo isOpen={isPopupOpenSubirArchivo} onClose={onClosePopupSubirArchivo} courseId={course.course_id} />
            <PopupSeeReview isOpen={isPopupOpenSeeReview} onClose={onClosePopupSeeReview} courseId={course.course_id} />
            <PopupDetail isOpen={isPopupOpenDetail} onClose={onClosePopupDetail} course={course} formattedDate={formattedDate} getProfesorName={getProfesorName} />

        </Card>
    );
};

export default Item;
