import PropTypes from 'prop-types';
import { Input, InputGroup, InputLeftElement, Box, Button } from '@chakra-ui/react';
import { SearchIcon } from '@chakra-ui/icons';
import React from "react";
import '../estilos/SearchBar.css'

const SearchBar = ({ onSearchResults }) => {
    const [searchTerm, setSearchTerm] = React.useState('');

    // Función para construir la URL con parámetros
    const buildSearchUrl = (baseUrl, searchTerm, limit = 20, offset = 1) => {
        const params = new URLSearchParams();
        params.append('limit', limit);
        params.append('offset', offset);

        if (searchTerm.trim() !== '') {
            params.append('q', searchTerm); // Agrega el término de búsqueda si no está vacío
        }

        return `${baseUrl}?${params.toString()}`;
    };

// Función principal para manejar la búsqueda
    const handleSearch = async (e) => {
        e.preventDefault(); // Evita el comportamiento por defecto del formulario

        const baseUrl = 'http://localhost:8082/search';
        const url = buildSearchUrl(baseUrl, searchTerm); // Construye la URL con los parámetros

        try {
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                onSearchResults(data); // Envía los resultados a la función de callback
            } else {
                handleNoResults(); // Maneja el caso de no encontrar resultados
            }
        } catch (error) {
            handleError(error); // Maneja errores de la solicitud
        }
    };

// Función para manejar la falta de resultados
    const handleNoResults = () => {
        alert("No se encontraron cursos.");
        onSearchResults([]); // Devuelve un arreglo vacío
    };

// Función para manejar errores
    const handleError = (error) => {
        console.error('Error al realizar la solicitud al backend:', error);
        alert("Error al buscar cursos. Inténtalo de nuevo más tarde.");
        onSearchResults([]); // Devuelve un arreglo vacío
    };


    return (
        <Box className='search' id='caja'>
            <form onSubmit={handleSearch}>
                <InputGroup>
                    <InputLeftElement pointerEvents="none">
                        <SearchIcon id='icono' />
                    </InputLeftElement>
                    <Input
                        className='input'
                        type="text"
                        placeholder="Buscar cursos por nombre..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </InputGroup>
                {/*<Button type="submit" mt={2} width="100%">Buscar</Button>*/}
            </form>
        </Box>
    );
};

SearchBar.propTypes = {
    onSearchResults: PropTypes.func.isRequired,
};

export default SearchBar;