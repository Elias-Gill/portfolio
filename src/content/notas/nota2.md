---
titulo: Notas de nada que ver
descripcion: Nada de nada, es solo para probar
---

# TODO:
- [X] Configurar routers
- [X] Disenar y armar la BD
- [ ] Crear Tests para los endpoints
- [ ] Crear un server de testing
- [ ] Levantar server

# Terminado
- [X] Actualizar las descripciones de swagger
- [X] Disenar las tareas y la estructura del guardado y parseo de los archivos excel. Tal vez un parseo automatico a json, guardando cada carrera (sheet) en un archivo separado.

# Diseno de la DB
La DB cuenta con 4 tablas (el arhivo de creacion ya se encuentra en el proyecto):
### usuarios
```sql
id_usuario serial4 NOT NULL,
nombre varchar NOT NULL,
email varchar NOT NULL,
imagen varchar NOT NULL, -- se usando conversion a string
contrasena varchar NOT NULL,
```

### hojas_excel
Contiene los nombres de los archivos descargados
```sql
id_excel serial NOT NULL,
fecha_carga date NOT NULL,
nombre_archivo varchar NOT NULL,
```
### carreras
Contiene las carreras disponibles dentro de un archivo excel (dios sabra si estoy haciendo bien las cosas).
```sql
id_carrera serial not null,
nombre varchar not null,
archivo_excel integer not null,
n_hoja int not null, -- numero de hoja en el excel
```

### materias: 
Contiene la informacion de todas las materias disponibles.
```sql
id_materia serial4 NOT NULL,
descripcion varchar NOT NULL,
archivo_excel integer NOT NULL, (FK)
carrera integer NOT NULL,       (FK)

lunes varchar,
martes varchar,
miercoles varchar,
jueves varchar,
viernes varchar,

parcial1 varchar,
parcial2 varchar,
final1 varchar,
final2 varchar,
```
### horarios
Contiene los horarios de los usuarios. Cada usuario puede contar con varios horarios.
```sql
id_horario serial4 NOT NULL,
id_usuario integer NOT NULL, (FK)
```

### detalle_horario
Tabla de enlace entre los horarios y las materias.
```sql
id_horario integer NOT NULL, (FK)
id_materia integer NOT NULL, (FK)
```

# Routers
## /users
- [X] **GET** /login Login de usuario
    - [X] **GET** /info Traer datos del usuario
    - [X] **POST** / Crear nuevo usuario
    - **PUT** / Modificar datos de un usuario
    - **DELETE** / Dar de baja un usuario

## /schedule
- **GET** / Traer los horarios del usuario
    - **POST** / Generar un nuevo horario para el usuario
    - **PUT** / Modificar datos de un horario
    - **DELETE** / Eliminar un horario
        - _/sheets_
            - **GET** _/_ Traer la lista de archivos excel disponible
            - **GET** _/subjects_ Traer la lista de materias contenidas en un archivo

# Tareas secundarias
- [ ] Aprender a hacer una imagen docker de esto

# Notas y pensamientos
REPLANTEANDOME PASAR TODO A MONGODB para hacer mas facil el guardado de la informacion de los archivos excel, ya que ya tengo el parser que me transforma luego todo a json.

O en todo caso conectar con una instancia de Google Drive o algo loco por ahi. O cambiar para que guarde los parseos como archivos temporales.
