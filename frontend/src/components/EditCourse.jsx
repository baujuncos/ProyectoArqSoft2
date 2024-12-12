import React, { useState, useEffect } from 'react';
import { Button, FormControl, FormLabel, Input } from "@chakra-ui/react";
import Cookies from 'js-cookie';

const EditCourse = ({ courseId, onClose }) => {
    const [courseData, setCourseData] = useState({
        nombre: '',
        categoria: '',
        descripcion: '',
        duracion: '',
        requisitos: '',
        url_image: '',
        fecha_inicio: '',
    });

    // Fetch para obtener los datos del curso actual
    const fetchCourse = async () => {
        try {
            const response = await fetch(`http://localhost:8081/courses/${courseId}`);
            if (response.ok) {
                const data = await response.json();
                setCourseData({
                    nombre: data.nombre,
                    categoria: data.categoria,
                    descripcion: data.descripcion,
                    duracion: data.duracion,
                    requisitos: data.requisitos,
                    url_image: data.url_image,
                    fecha_inicio: data.fecha_inicio.split('T')[0], // Formato YYYY-MM-DD
                });
            } else {
                console.error("Error al obtener los datos del curso.");
            }
        } catch (error) {
            console.error("Error al cargar los detalles del curso:", error);
        }
    };

    useEffect(() => {
        fetchCourse();
    }, [courseId]);

    // Maneja los cambios en los inputs
    const handleChange = (e) => {
        const { name, value } = e.target;
        setCourseData({ ...courseData, [name]: value });
    };

    // Envía los datos actualizados al backend
    const handleSubmit = async (e) => {
        e.preventDefault();

        // Formatear datos antes de enviarlos
        const formattedData = {
            ...courseData,
            duracion: parseInt(courseData.duracion, 10), // Convertir duración a número
            fecha_inicio: new Date(courseData.fecha_inicio).toISOString() // Convertir a formato ISO
        };

        console.log("Datos enviados al backend:", formattedData);

        try {
            const response = await fetch(`http://localhost:8081/courses/${courseId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formattedData), // Enviar formattedData
            });

            if (response.ok) {
                alert("Curso actualizado correctamente");
                onClose();
                window.location.reload();
            } else {
                const errorData = await response.json();
                console.error("Error al actualizar el curso:", errorData);
                alert(`Error al actualizar el curso: ${errorData.message}`);
            }
        } catch (error) {
            console.error("Error al actualizar el curso:", error);
            alert("Error de red al actualizar el curso.");
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <FormControl>
                <FormLabel>Nombre</FormLabel>
                <Input name="nombre" value={courseData.nombre} onChange={handleChange} />
            </FormControl>
            <FormControl>
                <FormLabel>Categoría</FormLabel>
                <Input name="categoria" value={courseData.categoria} onChange={handleChange} />
            </FormControl>
            <FormControl>
                <FormLabel>Descripción</FormLabel>
                <Input name="descripcion" value={courseData.descripcion} onChange={handleChange} />
            </FormControl>
            <FormControl>
                <FormLabel>Duración</FormLabel>
                <Input name="duracion" value={courseData.duracion} onChange={handleChange} />
            </FormControl>
            <FormControl>
                <FormLabel>Requisitos</FormLabel>
                <Input name="requisitos" value={courseData.requisitos} onChange={handleChange} />
            </FormControl>
            <FormControl>
                <FormLabel>URL de la imagen</FormLabel>
                <Input name="url_image" value={courseData.url_image} onChange={handleChange} />
            </FormControl>
            <FormControl>
                <FormLabel>Fecha de inicio</FormLabel>
                <Input type="date" name="fecha_inicio" value={courseData.fecha_inicio} onChange={handleChange} />
            </FormControl>
            <Button mt="4" colorScheme="blue" type="submit">
                Guardar cambios
            </Button>
        </form>
    );
};

export default EditCourse;
