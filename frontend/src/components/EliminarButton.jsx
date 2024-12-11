import React, { useEffect, useState } from 'react';
import Cookies from 'js-cookie';
import '../estilos/Inscribirmebutton.css';

const DeleteCourse = ({ courseId, onClose }) => {
    const tokenUser = Cookies.get('token');
    const [user_id, setUserId] = useState(null);
    const [isAdmin, setIsAdmin] = useState(false);

    useEffect(() => {
        const storedUserId = Cookies.get('user_id');
        if (storedUserId) {
            setUserId(parseInt(storedUserId, 10));
        }

        const storedAdmin = Cookies.get('admin');
        if (storedAdmin) {
            setIsAdmin(storedAdmin === "1"); // Si 'admin' es "1", entonces es admin
        }
    }, []);

    const handleDelete = async () => {
        try {
            const response = await fetch(`http://localhost:8081/courses/${courseId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    //'Authorization': `Bearer ${tokenUser}`
                }
            });

            if (response.ok) {
                alert('Curso eliminado exitosamente');
                window.location.reload();
            } else {
                const errorData = await response.json();
                alert(`Error al eliminar el curso: ${errorData.message}`);
            }
        } catch (error) {
            console.error(`Error de red al eliminar el curso: ${error.message}`);
            alert("Error al eliminar el curso");
        }
    };

    // Si no es admin, no mostrar el botón
    if (!user_id || !isAdmin) {
        return null;
    }

    return (
        <button className="delete-button" onClick={handleDelete}>ELIMINAR</button>
    );
};

export default DeleteCourse;
