# Lumen PHP Framework (Dibu Backend Proyect) 

https://bienestar.unas.edu.pe/home/services

## Instanciar entorno de desarrollo local

1.  Clonar y abrir el proyecto con vs-code.
    **CMD**

```
git clone https://github.com/iomarigor/dibu-backend.git
cd dibu-backend
code .
```

2. Configurar el .env
   **.env**

```
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=homestead
DB_USERNAME=homestead
DB_PASSWORD=secret
```

3. Acceder a la carpeta public y instalar las dependencias del proyecto

```
cd public
composer install
cd ..
```

4. Iniciar los contenedores de docker

```
docker-compose up -d
```

5. Acceder al proyecto por la url:
   http://localhost:8080/public/
   [![](https://i.ibb.co/r60d4Pn/Captura-de-pantalla-2023-07-22-231112.png)](https://i.ibb.co/r60d4Pn/Captura-de-pantalla-2023-07-22-231112.png)
## Otros comandos
```
php artisan make:migration *nombre_migracion* --create=*nombre_tabla* *//Crear las migraciones*
php artisan migrate                                                   *//Correr las migraciones*
php artisan make:controller DetalleVentasController                   *//Crear controlador sin recursos*
php artisan make:controller VentasController --resource               *//Crear controlador con recursos
```

## Colaboradores

- Daniel Erick Ruiz Balcazar (daniel.ruiz@unas.edu.pe).
- Carlos Andres Cuadros Cordoba (carlos.cuadros@unas.edu.pe).
- Luis Lucero Saldaña (luis.lucero@unas.edu.pe).
- Jordi Brandon Valverde Chavez (jordi.valverde@unas.edu.pe).
- Diego Valencia Moya (diego.valencia@unas.edu.pe).
- Iomar Igor Alegre Barrera (iomar.alegre@unas.edu.pe).

## Colaboradores 2
- Samaniego Martel Alvaro (alvaro.samaniego@unas.edu.pe)
- Rojas Rivadeneiro Yonil (yonil.rojas@unas.edu.pe)
