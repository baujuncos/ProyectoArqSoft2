import React from 'react';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalCloseButton,
    ModalBody,
    Text,
    Button } from "@chakra-ui/react";

const PopupDetail = ({ isOpen, onClose, course }) => {
    return (
        <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Ver detalle</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <Text><strong>Nombre del curso:</strong> {course.name}</Text>
                    <Text><strong>Descripción:</strong> {course.descripcion}</Text>
                    <Text><strong>Requisitos:</strong> {course.requisitos}</Text>
                    <Text><strong>Duración:</strong> {course.duracion}</Text>
                    <Button onClick={onClose}>Cerrar</Button>
                </ModalBody>
            </ModalContent>
        </Modal>
    );
};

export default CourseDetailPopup;