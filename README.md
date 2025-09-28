# Setup

- `BLOG_PATH`:
  path a la carpeta donde esta contenido el contenido del repositorio para el blog.
- `WEBHOOK_SECRET`:
  la key para conectar al webhook de actualizaciones del repositorio de posts.

Correr el proyecto.

## Notas sobre como funciona

El blog se sirve de manera dinamica, usando directametne el nombre de archivo (se usa
sanitizacion del nombre de archivo antes de tratar de servir el archivo).
Se extraen los metadatos y se parsea usnado `goldmark`.

Existen 2 `sitemap.xml`, uno el cual se encuentra en este repositorio, el segundo se encuentra
dentro del repositorio de blogs, ahi se puede actualizar de manera manual con cada nuevo post.

Cuando se hace un push al repo de blogs, el webhook levanta al server en un endpoint especial
para avisarle que actualice su carpeta de blogs con un git pull, ademas de hacer ping a Google
para actualizar el sitemap del blog.

Las imagenes se encuentran desde la misma carpeta del blog dentro de `media/`.
Para hacer referencia a imagenes ahi, simplemente usar el url `/media/...`, el server
automaticamente sirve desde ese media path.

Licensed MIT.
