START TRANSACTION;

-- Crear la tabla `users`
CREATE TABLE IF NOT EXISTS users (
                                     user_id INT AUTO_INCREMENT PRIMARY KEY,
                                     email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    admin BOOLEAN NOT NULL
    );

-- Insertar datos en `users`
INSERT INTO users (email, password, nombre, apellido, admin) VALUES
                                                                 ('pauliortiz@example.com', MD5('contraseña1'), 'paulina', 'ortiz', TRUE),
                                                                 ('baujuncos@example.com', MD5('contraseña2'), 'bautista', 'juncos', TRUE),
                                                                 ('belutreachi2@example.com', MD5('contraseña3'), 'belen', 'treachi', FALSE),
                                                                 ('virchurodiguez@example.com', MD5('contraseña4'), 'virginia', 'rodriguez', FALSE),
                                                                 ('johndoe@example.com', MD5('contraseña5'), 'John', 'Doe', FALSE),
                                                                 ('alicesmith@example.com', MD5('contraseña6'), 'Alice', 'Smith', TRUE),
                                                                 ('bobjohnson@example.com', MD5('contraseña7'), 'Bob', 'Johnson', FALSE),
                                                                 ('janedoe@example.com', MD5('contraseña8'), 'Jane', 'Doe', FALSE),
                                                                 ('emilywilliams@example.com', MD5('contraseña9'), 'Emily', 'Williams', TRUE);

-- Crear la tabla `inscripciones`
CREATE TABLE IF NOT EXISTS inscripciones (
                                             inscripcion_id INT AUTO_INCREMENT PRIMARY KEY,
                                             id_usuario INT NOT NULL,
                                             id_curso VARCHAR(50) NOT NULL,
    fecha_inscripcion DATETIME NOT NULL,
    FOREIGN KEY (id_usuario) REFERENCES users(user_id)
    );

-- Insertar datos en `inscripciones`
INSERT INTO inscripciones (id_usuario, id_curso, fecha_inscripcion) VALUES
                                                                        (3, '1', NOW()),
                                                                        (5, '2', NOW()),
                                                                        (6, '3', NOW());

COMMIT;