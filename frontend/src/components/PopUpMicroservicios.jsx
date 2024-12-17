import { useState } from "react";
import {
    Drawer,
    DrawerOverlay,
    DrawerContent,
    DrawerCloseButton,
    DrawerHeader,
    DrawerBody,
    Button,
    VStack,
    Box,
    Spinner,
} from "@chakra-ui/react";
import Cookies from "js-cookie";

const MicroserviciosPopup = ({ isOpen, onClose }) => {
    const [containers, setContainers] = useState([]); // Lista de contenedores
    const [loadingContainers, setLoadingContainers] = useState(false); // Estado de carga
    const [errorContainers, setErrorContainers] = useState(null); // Estado de error

    const token = Cookies.get("token");

    // Fetch de los datos de contenedores
    const fetchContainers = async () => {
        setLoadingContainers(true);
        setErrorContainers(null); // Limpiar errores previos

        if (!token) {
            setErrorContainers("Token no encontrado. Por favor inicia sesión.");
            setLoadingContainers(false);
            return;
        }

        try {
            const response = await fetch("http://localhost:8004/microservices", {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || "Error al obtener los contenedores.");
            }

            const containerData = await response.json();
            setContainers(containerData.services || []); // Aseguramos la estructura
        } catch (error) {
            console.error("Error fetching containers:", error);
            setErrorContainers("Error al obtener la lista de contenedores.");
        } finally {
            setLoadingContainers(false);
        }
    };

    return (
        <Drawer isOpen={isOpen} placement="right" onClose={onClose}>
            <DrawerOverlay />
            <DrawerContent>
                <DrawerCloseButton />
                <DrawerHeader>Gestión de Microservicios</DrawerHeader>
                <DrawerBody>
                    <VStack spacing={4}>
                        <Button
                            w="100%"
                            onClick={fetchContainers}
                            isLoading={loadingContainers}
                            colorScheme="blue"
                            style={{ fontFamily: "Spoof Trial, sans-serif" }}
                        >
                            Ver contenedores
                        </Button>

                        {loadingContainers && <Spinner size="xl" />}

                        {errorContainers && (
                            <Box w="100%" color="red.500">
                                {errorContainers}
                            </Box>
                        )}

                        {!loadingContainers && containers.length > 0 && (
                            <Box w="100%" className="containerGrid" p={2}>
                                {containers.map((service, index) => (
                                    <Box
                                        key={index}
                                        className="containerCard"
                                        p={4}
                                        border="1px solid #ccc"
                                        borderRadius="md"
                                        boxShadow="sm"
                                    >
                                        <p>
                                            <strong>Nombre:</strong> {service.name}
                                        </p>
                                        <p>
                                            <strong>ID del Contenedor:</strong> {service.container}
                                        </p>
                                        <p>
                                            <strong>Puerto:</strong> {service.port}
                                        </p>
                                        <p>
                                            <strong>Estado:</strong> {service.state}
                                        </p>
                                    </Box>
                                ))}
                            </Box>
                        )}
                        {!loadingContainers && containers.length === 0 && !errorContainers && (
                            <Box w="100%" color="gray.500" textAlign="center">
                                No hay microservicios disponibles.
                            </Box>
                        )}
                    </VStack>
                </DrawerBody>
            </DrawerContent>
        </Drawer>
    );
};

export default MicroserviciosPopup;
