import React from 'react';
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

const PopupDetail = ({ isOpen, onClose, course, formattedDate, getProfesorName }) => {
    return (
        <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Ver detalle</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <Text py='2' className="card-text">Categoría: {course.categoria}</Text>
                    <Text className="card-textt">Duracion: {course.duracion}hs</Text>
                    <Text className="card-textt">Fecha de inicio: {formattedDate}</Text>
                    <Text className="card-textt">Requisito: Nivel {course.requisitos}</Text>
                    <Text className="card-textt">Profesor: {getProfesorName(course.profesor_id)}</Text>
                    <Button onClick={onClose}>Cerrar</Button>
                    <Inscribirmebutton courseId={course.course_id} />
                </ModalBody>
            </ModalContent>
        </Modal>
    );
};

export default PopupDetail;
