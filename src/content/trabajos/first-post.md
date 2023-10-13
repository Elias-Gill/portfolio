---
titulo: 'First post'
descripcion: 'Lorem ipsum dolor sit amet'
pubDate: 'Jul 08 2022'
caption: '/placeholder-hero.jpg'
---

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vitae ultricies leo integer malesuada nunc vel risus commodo viverra. Adipiscing enim eu turpis egestas pretium. Euismod elementum nisi quis eleifend quam adipiscing. In hac habitasse platea dictumst vestibulum. Sagittis purus sit amet volutpat. Netus et malesuada fames ac turpis egestas. Eget magna fermentum iaculis eu non diam phasellus vestibulum lorem. Varius sit amet mattis vulputate enim. Habitasse platea dictumst quisque sagittis. Integer quis auctor elit sed vulputate mi. Dictumst quisque sagittis purus sit amet.

Morbi tristique senectus et netus. Id semper risus in hendrerit gravida rutrum quisque non tellus. Habitasse platea dictumst quisque sagittis purus sit amet. Tellus molestie nunc non blandit massa. Cursus vitae congue mauris rhoncus. Accumsan tortor posuere ac ut. Fringilla urna porttitor rhoncus dolor. Elit ullamcorper dignissim cras tincidunt lobortis. In cursus turpis massa tincidunt dui ut ornare lectus. Integer feugiat scelerisque varius morbi enim nunc. Bibendum neque egestas congue quisque egestas diam. Cras ornare arcu dui vivamus arcu felis bibendum. Dignissim suspendisse in est ante in nibh mauris. Sed tempus urna et pharetra pharetra massa massa ultricies mi.

Mollis nunc sed id semper risus in. Convallis a cras semper auctor neque. Diam sit amet nisl suscipit. Lacus viverra vitae congue eu consequat ac felis donec. Egestas integer eget aliquet nibh praesent tristique magna sit amet. Eget magna fermentum iaculis eu non diam. In vitae turpis massa sed elementum. Tristique et egestas quis ipsum suspendisse ultrices. Eget lorem dolor sed viverra ipsum. Vel turpis nunc eget lorem dolor sed viverra. Posuere ac ut consequat semper viverra nam. Laoreet suspendisse interdum consectetur libero id faucibus. Diam phasellus vestibulum lorem sed risus ultricies tristique. Rhoncus dolor purus non enim praesent elementum facilisis. Ultrices tincidunt arcu non sodales neque. Tempus egestas sed sed risus pretium quam vulputate. Viverra suspendisse potenti nullam ac tortor vitae purus faucibus ornare. Fringilla urna porttitor rhoncus dolor purus non. Amet dictum sit amet justo donec enim.

Mattis ullamcorper velit sed ullamcorper morbi tincidunt. Tortor posuere ac ut consequat semper viverra. Tellus mauris a diam maecenas sed enim ut sem viverra. Venenatis urna cursus eget nunc scelerisque viverra mauris in. Arcu ac tortor dignissim convallis aenean et tortor at. Curabitur gravida arcu ac tortor dignissim convallis aenean et tortor. Egestas tellus rutrum tellus pellentesque eu. Fusce ut placerat orci nulla pellentesque dignissim enim sit amet. Ut enim blandit volutpat maecenas volutpat blandit aliquam etiam. Id donec ultrices tincidunt arcu. Id cursus metus aliquam eleifend mi.

Tempus quam pellentesque nec nam aliquam sem. Risus at ultrices mi tempus imperdiet. Id porta nibh venenatis cras sed felis eget velit. Ipsum a arcu cursus vitae. Facilisis magna etiam tempor orci eu lobortis elementum. Tincidunt dui ut ornare lectus sit. Quisque non tellus orci ac. Blandit libero volutpat sed cras. Nec tincidunt praesent semper feugiat nibh sed pulvinar proin gravida. Egestas integer eget aliquet nibh praesent tristique magna.

# Fonts (como llamarlas desde CSS)

INFO: los titulos tambien son sus nombres de clase CSS
titulo - font-family: 'Secular One', sans-serif;
texto - font-family: 'Murecho', sans-serif;

# Tailwind

- Para cambiar de fuente: font-<fuente>
- Para padding: px-... o py-...

# Astrojs

## Estructura de carpetas en Astro

- En la carpeta _index_ se encuentran todas las paginas que se van a mostrar y astro renderiza de a uno
- Entiendo que para las paginas que va a tener tu blog metes las paginas dentro de _pages_. Podes poner paginas individuales o podes crear carpetas que son como los _endpoints de una api_. Dentro de esa carpeta Astro busca por el archivo _index.astro_. Para hacer Referencias a estas paginas agregas el nombre del archivo o carpeta en el index (funciona recursivamente, como los _endpoints_ de una API).
- _BaseHead_ parece ser la metadata del archivo, mejor no tocar la config original de astro.
- El archivo _[...slug].astro_ se utiliza como comodin (o plantilla) en funcion al url que se proporciona. En este proyecto se esta trayendo todos el contenido dentro de 'content/blog/' utilizando la funcion:

- Dentro de la carpeta donde se crea una nueva coleccion ("content") se debe poner en config.ts la configuracion de la coleccion:

### Layout vs components
Los componentes son nada de 

## Sintaxis de astro

- En astro todo lo que esta entre tres guiones (---) es <script></script>

- Los layouts sirven para poder hacer eso, marcos de aplicacion pero cuando usas markdown supongo, son invocados automáticamente por AstroJS. Se encuentran en la carpeta 'layouts'. Dentro de los layouts
puedo tener diferentes elementos, de los cuales el mas importante seguro es _<slot />_ el cual se utiliza para inyectar html
externo al momento de utilizar el layout.

```astro
---
import type { HTMLAttributes } from 'astro/types';

type Props = HTMLAttributes<'a'>;

const { href, class: className, ...props } = Astro.props;

const { pathname } = Astro.url;
const isActive = href === pathname || href === pathname.replace(/\/$/, '');
---

<a href={href} class:list={[className, { active: isActive }]} {...props}>
    <slot />
</a>
<style>
a {
    display: inline-block;
    text-decoration: none;
}
a.active {
    font-weight: bolder;
    text-decoration: underline;
}
</style>
```
## Mas importante de Astro

- _RSS_: Astro proporciona una generación rápida y automática de RSS feeds para blogs u otros sitios web con mucho contenido. Los feeds RSS proporcionan una forma fácil para que los usuarios se suscriban a tu contenido. Un feed RSS es una forma de distribuir contenido actualizado a los suscriptores de un sitio web. Permite a los usuarios recibir notificaciones automáticas cuando se publica nuevo contenido en el sitio.
- Los layouts definen la estructura general de una página y se utilizan para definir cómo se organizan los componentes en una página. Pueden ser utilizados en varias páginas que siguen la misma estructura general, pero con contenido diferente.
- Los componentes son bloques de construcción reutilizables que se utilizan para construir la UI. Pueden ser utilizados dentro de otros componentes o layouts para construir una UI más avanzada.

### Ejemplo config.ts de content
```ts
import { defineCollection, z } from 'astro:content';

const works = defineCollection({
    // Type-check frontmatter using a schema
    schema: z.object({
        title: z.string(),
        description: z.string(),
        // Traesform string to Date object
        pubDate: z.string(),
        caption: z.string().optional(),
    }),
});

export const collections = { trabajos };
```

### Usage of a collection entry

```astro
---
// para plantillas
type Props = CollectionEntry<'notas'>['data'];
const { titulo, descripcion } = Astro.props;

// usar dentro del index
const notas = (await getCollection('notas'))
---
```

### ...slug.astro
```astro
---
import { CollectionEntry, getCollection } from 'astro:content';
import Note from '../../layouts/Note.astro';

export async function getStaticPaths() {
    const notas = await getCollection('notas');
    return notas.map((nota) => ({
        params: { slug: nota.slug },
        props: nota,
    }));
}
type Props = CollectionEntry<'notas'>;

const post = Astro.props;
const { Content } = await post.render();
---

<Note {...post.data}>
    <h1>{post.data.titulo}</h1>
    <Content />
</Note>
```
