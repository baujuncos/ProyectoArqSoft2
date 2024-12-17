import React from 'react';
import Cookies from 'js-cookie';
import { Button, Input, InputGroup, InputRightElement } from "@chakra-ui/react";
import '../estilos/Login.css';

const PasswordInput = ({ password, setPassword }) => {
    const [show, setShow] = React.useState(false);
    const handleClick = () => setShow(!show);

    return (
        <InputGroup size='md' className='inputContrasena'>
            <Input
                pr='4.5rem'
                type={show ? 'text' : 'password'}
                placeholder='Ingrese su contraseña...'
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                style={{ border: '2px solid black', fontFamily: 'Spoof Trial, sans-serif' }}
            />
            <InputRightElement width='4.5rem'>
                <Button h='1.75rem' size='sm' onClick={handleClick} style={{ fontFamily: 'Spoof Trial, sans-serif' }}>
                    {show ? 'Ocultar' : 'Mostrar'}
                </Button>
            </InputRightElement>
        </InputGroup>
    );
};

const Login = ({ onClose }) => {
    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');

    Cookies.set('user_id', 0);

    const handleSubmit = async (e) => {
        e.preventDefault(); // Evita que se recargue la página
        let emptyLogin = false;

        if (email === '') {
            alert('El campo de correo electrónico es obligatorio.');
            emptyLogin = true;
        }

        if (password === '') {
            alert('El campo de contraseña es obligatorio.');
            emptyLogin = true;
        }

        if (!emptyLogin) {
            try {
                const response = await fetch('http://localhost:8080/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password }),
                });

                if (response.ok) {
                    const data = await response.json();
                    console.log('Respuesta del servidor:', data);

                    if (data.user_id) {
                        Cookies.set('user_id', data.user_id, { expires: 1, path: '/' });
                        Cookies.set('email', email, { expires: 1, path: '/' });
                        Cookies.set('token', data.token, { expires: 1, path: '/' });
                        Cookies.set('admin', data.admin ? "1" : "0", { expires: 1, path: '/' });

                        onClose(); // Cierra el popup solo si todo es exitoso
                        window.location.reload();
                    }
                } else {
                    alert("Usuario no registrado o contraseña incorrecta.");
                }
            } catch (error) {
                Cookies.set('user_id', "-1");
                console.log('Error al realizar la solicitud al backend:', error);
            }
        }
    };

    return (
        <form id="formLogin" onSubmit={handleSubmit}>
            <Input
                id={'inputEmailLogin'}
                placeholder="Ingrese su correo electrónico..."
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                style={{ border: '2px solid black', fontFamily: 'Spoof Trial, sans-serif' }}
            />
            <PasswordInput password={password} setPassword={setPassword} />
            <Button colorScheme="blue" mr={3} type="submit">
                Iniciar sesión
            </Button>
        </form>
    );
};

export default Login;
