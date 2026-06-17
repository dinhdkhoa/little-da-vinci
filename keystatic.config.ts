import { config, fields, singleton, collection } from '@keystatic/core';

// Reusable localized fields
const localeString = (label: string) => fields.object({
  en: fields.text({ label: `${label} (English)` }),
  vi: fields.text({ label: `${label} (Tiếng Việt)` }),
});

const localeText = (label: string, multiline = true) => fields.object({
  en: fields.text({ label: `${label} (English)`, multiline }),
  vi: fields.text({ label: `${label} (Tiếng Việt)`, multiline }),
});

// Cloudinary Asset fallback
const cloudinaryField = (label: string) => fields.object({
  secure_url: fields.url({ label: `${label} - Cloudinary URL` }),
  public_id: fields.text({ label: 'Public ID (Optional)' }),
});

export default config({
  storage: process.env.NODE_ENV === 'development' ? {
    kind: 'local',
  } : {
    kind: 'github',
    repo: 'dinhdkhoa/little-da-vinci',
  },
  singletons: {
    settings: singleton({
      label: 'Site Settings',
      path: 'src/content/settings',
      format: { data: 'json' },
      schema: {
        siteTitle: localeString('Site Title'),
        siteDescription: localeText('Site Description'),
        footerTagline: localeString('Footer Tagline'),
        navigation: fields.object({
          home: localeString('Home (Brand Name)'),
          about: localeString('About'),
          services: localeString('Services'),
          eventLab: localeString('The Event Lab'),
          contact: localeString('Contact'),
          bookCall: localeString('Book Call Button'),
          blog: localeString('Blog'),
        }),
        blogLabels: fields.object({
          title: localeString('Blog Title'),
          subtitle: localeString('Blog Subtitle'),
          readMore: localeString('Read More Label'),
          backToList: localeString('Back to List Label'),
          noPostsFound: localeString('No Posts Found Message'),
          publishedOn: localeString('Published On Label'),
          relatedPosts: localeString('Related Posts Title'),
        })
      }
    }),
    home: singleton({
      label: 'Home Page',
      path: 'src/content/home',
      format: { data: 'json' },
      schema: {
        hero: fields.object({
          tagline: localeString('Tagline'),
          title: localeText('Main Title (HTML tags allowed)'),
          description: localeText('Description'),
          cta1_label: localeString('Primary CTA Label'),
          cta1_link: fields.text({ label: 'Primary CTA Link' }),
          cta2_label: localeString('Secondary CTA Label'),
          cta2_link: fields.text({ label: 'Secondary CTA Link' }),
          image: cloudinaryField('Hero Image'),
        }),
        howWeTeach: fields.object({
          title: localeString('Section Title'),
          description: localeText('Description'),
          steps: fields.array(
            fields.object({
              icon: fields.text({ label: 'Material Symbol Icon Name' }),
              title: localeString('Title'),
              description: localeText('Description'),
            }),
            { label: 'Steps', itemLabel: props => 'Step' }
          )
        }),
        programs: fields.object({
          title: localeString('Section Title'),
          quote: localeText('Quote'),
          courses: fields.array(
            fields.object({
              image: cloudinaryField('Course Image'),
              ageRange: localeString('Age Range'),
              title: localeString('Course Title'),
              description: localeText('Course Description'),
              tags: fields.array(
                localeString('Tag'),
                { label: 'Tags', itemLabel: props => 'Tag' }
              ),
            }),
            { label: 'Courses', itemLabel: props => 'Course' }
          )
        }),
        whyChooseUs: fields.object({
          title: localeString('Section Title'),
          features: fields.array(
            fields.object({
              icon: fields.text({ label: 'Material Symbol Icon Name' }),
              title: localeString('Feature Title'),
              description: localeText('Feature Description'),
            }),
            { label: 'Features', itemLabel: props => 'Feature' }
          )
        }),
        gallery: fields.object({
          title: localeString('Section Title'),
          images: fields.array(
            cloudinaryField('Gallery Image'),
            { label: 'Images', itemLabel: props => 'Image' }
          )
        }),
        testimonials: fields.object({
          items: fields.array(
            fields.object({
              content: localeText('Testimonial Content'),
              author: localeString('Author Name'),
            }),
            { label: 'Testimonials', itemLabel: props => 'Testimonial' }
          )
        }),
        cta: fields.object({
          title: localeString('Section Title'),
          description: localeText('Description'),
          address: localeString('Address'),
          email: localeString('Email'),
          phone: localeString('Phone Number'),
        })
      }
    })
  },
  collections: {
    authors: collection({
      label: 'Authors',
      slugField: 'name',
      path: 'src/content/authors/*',
      format: { data: 'json' },
      schema: {
        name: fields.slug({ name: { label: 'Name' } }),
        image: cloudinaryField('Image'),
        bio: localeText('Bio')
      }
    }),
    categories: collection({
      label: 'Categories',
      slugField: 'title_internal',
      path: 'src/content/categories/*',
      format: { data: 'json' },
      schema: {
        title_internal: fields.slug({ name: { label: 'Internal Title' } }),
        title: localeString('Title'),
        description: localeText('Description')
      }
    }),
    posts: collection({
      label: 'Blog Posts',
      slugField: 'slug',
      path: 'src/content/posts/*',
      format: { data: 'json' },
      schema: {
        slug: fields.slug({ name: { label: 'Internal Slug' } }),
        title: localeString('Title'),
        slug_en: fields.text({ label: 'English Slug (Required for URL)' }),
        slug_vi: fields.text({ label: 'Vietnamese Slug (Required for URL)' }),
        author: fields.relationship({ label: 'Author', collection: 'authors' }),
        mainImage: cloudinaryField('Main Image'),
        mainImageAlt: fields.text({ label: 'Main Image Alt text' }),
        categories: fields.array(
          fields.relationship({ label: 'Category', collection: 'categories' }),
          { label: 'Categories', itemLabel: props => props.value || 'Category' }
        ),
        publishedAt: fields.date({ label: 'Published at' }),
        excerpt: localeText('Excerpt', true),
        body_en: fields.document({ label: 'Body (English)', formatting: true, dividers: true, links: true, images: true }),
        body_vi: fields.document({ label: 'Body (Tiếng Việt)', formatting: true, dividers: true, links: true, images: true }),
      }
    })
  }
});