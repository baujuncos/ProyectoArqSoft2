import React, { useEffect, useState } from 'react';
import Cookies from 'js-cookie';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalCloseButton,
    ModalBody,
    Text,
    Button
} from "@chakra-ui/react";
import Inscribirmebutton from './Inscribirmebutton.jsx';
import EliminarButton from "./EliminarButton.jsx";
import PopupEdit from "./PopUpEdit.jsx";

const PopupDetail = ({ isOpen, onClose, course, formattedDate, getProfesorName }) => {
    const [isAdmin, setIsAdmin] = useState(false);
    const [isEditOpen, setIsEditOpen] = useState(false);

    useEffect(() => {
        const storedAdmin = Cookies.get('admin');
        if (storedAdmin) {
            setIsAdmin(storedAdmin === "1");
        }
    }, []);

    const openEditModal = () => setIsEditOpen(true);
    const closeEditModal = () => setIsEditOpen(false);

    return (
        <>
            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Ver detalle</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Text py='2' className="card-text">Categoría: {course.categoria}</Text>
                        <Text className="card-textt">Duracion: {course.duracion}hs</Text>
                        <Text className="card-textt">Fecha de inicio: {formattedDate}</Text>
                        <Text className="card-textt">Requisito: {course.requisitos}</Text>
                        <Text className="card-textt">Profesor: {getProfesorName(course.profesor_id)}</Text>

                        {/* Botón de Inscribirse */}
                        <Inscribirmebutton courseId={course.course_id} />

                        {/* Botón de Eliminar para administradores */}
                        {isAdmin && (
                            <>
                                <EliminarButton courseId={course.course_id} />
                                <Button
                                    w="40%" style={{ fontFamily: 'Spoof Trial, sans-serif' }}
                                    onClick={openEditModal}
                                >
                                    Editar
                                </Button>
                            </>
                        )}
                    </ModalBody>
                </ModalContent>
            </Modal>

            {/* Popup de Editar */}
            <PopupEdit isOpen={isEditOpen} onClose={closeEditModal} courseId={course.course_id} />
        </>
    );
};

export default PopupDetail;
