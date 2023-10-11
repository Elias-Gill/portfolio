import { defineCollection, z } from 'astro:content';

const trabajos = defineCollection({
    // Type-check frontmatter using a schema
    schema: z.object({
        title: z.string(),
        description: z.string(),
        // Traesform string to Date object
        pubDate: z.string(),
        caption: z.string().optional(),
    }),
});

const notas = defineCollection({
    // Type-check frontmatter using a schema
    schema: z.object({
        titulo: z.string(),
        descripcion: z.string(),
    }),
});

export const collections = { trabajos, notas };
