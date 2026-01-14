// Global variables
let isUploading = false;

document.getElementById('page-form').addEventListener('submit', async function(e) {
    e.preventDefault();

    const title = document.getElementById('title').value.trim();
    const slug = document.getElementById('slug').value.trim();

    // Get content from EasyMDE if editor is initialized, otherwise from textarea
    let content;
    if (easyMDE) {
        content = easyMDE.value().trim();
    } else {
        content = document.getElementById('content').value.trim();
    }
    const showInNav = document.getElementById('showInNav').checked;
    const sortOrder = parseInt(document.getElementById('sortOrder').value) || 0;

    // Basic validation
    if (!title) {
        alert('Please enter a title.');
        return;
    }
    if (!slug) {
        alert('Please enter a slug.');
        return;
    }
    if (!content) {
        alert('Please enter content.');
        return;
    }

    // Validate slug format
    const slugRegex = /^[a-z0-9-]+$/;
    if (!slugRegex.test(slug)) {
        alert('Slug must contain only lowercase letters, numbers, and hyphens.');
        return;
    }

    const pageData = {
        title: title,
        slug: slug,
        content: content,
        show_in_nav: showInNav,
        sort_order: sortOrder
    };

    const isEdit = window.location.pathname.includes('/edit');
    const url = isEdit ? `/api/pages/${pageId}` : '/api/pages';
    const method = isEdit ? 'PUT' : 'POST';
    const successMessage = isEdit ? 'Page updated successfully!' : 'Page created successfully!';

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(pageData)
        });

        if (response.ok) {
            alert(successMessage);
            window.location.href = '/admin/pages';
        } else if (response.status === 409) {
            alert('Error: Slug already exists. Please choose a different slug.');
        } else {
            const errorData = await response.json();
            alert(errorData.error || `Error ${isEdit ? 'updating' : 'creating'} page. Please try again.`);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Network error. Please try again.');
    }
});

// Auto-slugify title
document.getElementById('title').addEventListener('input', function() {
    const title = this.value.trim();
    const slugField = document.getElementById('slug');
    if (slugField.value === '' || slugField.value === slugify(slugField.dataset.original || '')) {
        slugField.value = slugify(title);
        slugField.dataset.original = title;
    }
});

function slugify(text) {
    return text
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .replace(/[\s_-]+/g, '-')
        .replace(/^-+|-+$/g, '');
}

// Initialize EasyMDE editor
let easyMDE;
document.addEventListener('DOMContentLoaded', function() {
    const contentTextarea = document.getElementById('content');
    if (contentTextarea && typeof EasyMDE !== 'undefined') {
        easyMDE = new EasyMDE({
            element: contentTextarea,
            spellChecker: false,
            renderingConfig: {
                singleLineBreaks: false,
                codeSyntaxHighlighting: true,
            },
            toolbar: [
                'bold', 'italic', 'heading', '|',
                'code', 'quote', 'unordered-list', 'ordered-list', '|',
                'link', 'image', '|',
                'preview', 'side-by-side', 'fullscreen', '|',
                'guide'
            ],
            status: ['autosave', 'lines', 'words', 'cursor'],
            autofocus: false,
            placeholder: 'Write your page content here...',
            minHeight: '400px',
            uploadImage: true,
            imageUploadFunction: function(file, onSuccess, onError) {
                const formData = new FormData();
                formData.append('image', file);

                fetch('/api/upload/image', {
                    method: 'POST',
                    body: formData
                })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Upload failed');
                    }
                    return response.json();
                })
                .then(data => {
                    onSuccess(data.data.filePath);
                })
                .catch(error => {
                    console.error('Upload error:', error);
                    onError('Upload failed. Please try again.');
                });
            }
        });
    } else {
        console.warn('EasyMDE not loaded, using plain textarea');
    }

    // Check if editing after editor is initialized
    const pathMatch = window.location.pathname.match(/^\/admin\/pages\/(\d+)\/edit$/);
    if (pathMatch) {
        const pageId = pathMatch[1];
        // Fetch page data
        fetch(`/api/pages/${pageId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Page not found');
                }
                return response.json();
            })
            .then(page => {
                document.getElementById('title').value = page.title || '';
                document.getElementById('slug').value = page.slug || '';
                // Set content in EasyMDE if available, otherwise textarea
                if (easyMDE) {
                    easyMDE.value(page.content || '');
                } else {
                    document.getElementById('content').value = page.content || '';
                }
                document.getElementById('showInNav').checked = page.show_in_nav || false;
                document.getElementById('sortOrder').value = page.sort_order || 0;
                document.getElementById('slug').dataset.original = page.slug || '';

                // Update page title
                document.title = `Edit Page: ${page.title}`;
                const heading = document.querySelector('h2');
                if (heading) {
                    heading.textContent = `Edit Page: ${page.title}`;
                }
            })
            .catch(error => {
                console.error('Error loading page:', error);
                alert('Error loading page data.');
            });
    }
});