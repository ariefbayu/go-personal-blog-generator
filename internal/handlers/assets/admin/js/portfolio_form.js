// Global variables
let isUploading = false;
let galleryImages = [];

document.addEventListener('DOMContentLoaded', function () {
  const isEdit = window.location.pathname.includes('/edit');
  let portfolioId = null;
  const pathMatch = window.location.pathname.match(/\/portfolio\/(\d+)\/edit$/);
  if (pathMatch) {
    portfolioId = pathMatch[1];
    loadPortfolioItem(portfolioId);
  }

  // Auto-slugify title (only when creating)
  const titleField = document.getElementById('title');
  const slugField = document.getElementById('slug');
  if (titleField && slugField) {
    titleField.addEventListener('input', function () {
      if (!isEdit) {
        const title = this.value.trim();
        if (slugField.value === '' || slugField.value === slugify(slugField.dataset.original || '')) {
          slugField.value = slugify(title);
          slugField.dataset.original = title;
        }
      }
    });
  }

  // Set up showcase image upload handler
  const imageUpload = document.getElementById('showcaseImageUpload');
  const imageUrlInput = document.getElementById('showcaseImageURL');
  const imagePreview = document.getElementById('imagePreview');
  const previewImg = document.getElementById('previewImg');

  if (imageUpload) {
    imageUpload.addEventListener('change', function (e) {
      const file = e.target.files[0];
      if (file) {
        if (!file.type.startsWith('image/')) {
          alert('Please select an image file');
          return;
        }

        if (file.size > 5 * 1024 * 1024) {
          alert('File size must be less than 5MB');
          return;
        }

        const slug = slugField.value.trim();
        if (!slug) {
          alert('Please enter a title or slug first so images can be organized properly.');
          imageUpload.value = '';
          return;
        }

        isUploading = true;
        imageUpload.disabled = true;

        const formData = new FormData();
        formData.append('image', file);

        fetch(`${window.ROOT_PREFIX || ''}/api/upload/image?slug=${encodeURIComponent(slug)}`, {
          method: 'POST',
          body: formData
        })
          .then(response => response.json())
          .then(data => {
            if (data.data && data.data.filePath) {
              imageUrlInput.value = data.data.filePath;
              previewImg.src = data.data.filePath;
              imagePreview.classList.remove('hidden');
            } else {
              alert('Upload failed: Invalid response format');
              console.error('Upload response:', data);
            }
          })
          .catch(error => {
            console.error('Upload error:', error);
            alert('Upload failed. Please try again.');
          })
          .finally(() => {
            isUploading = false;
            imageUpload.disabled = false;
          });
      }
    });
  }

  // Set up gallery images upload handler
  const galleryUpload = document.getElementById('portfolioImagesUpload');
  if (galleryUpload) {
    galleryUpload.addEventListener('change', function (e) {
      const files = Array.from(e.target.files);
      if (files.length === 0) return;

      const slug = slugField.value.trim();
      if (!slug) {
        alert('Please enter a title or slug first so gallery images can be organized properly.');
        galleryUpload.value = '';
        return;
      }

      // Upload files sequentially or in parallel
      isUploading = true;
      galleryUpload.disabled = true;

      let uploadPromises = files.map(file => {
        if (!file.type.startsWith('image/')) {
          return Promise.reject(new Error(`File ${file.name} is not an image`));
        }
        if (file.size > 5 * 1024 * 1024) {
          return Promise.reject(new Error(`File ${file.name} is larger than 5MB`));
        }

        const formData = new FormData();
        formData.append('image', file);

        return fetch(`${window.ROOT_PREFIX || ''}/api/upload/image?slug=${encodeURIComponent(slug)}`, {
          method: 'POST',
          body: formData
        })
          .then(res => {
            if (!res.ok) throw new Error(`Failed to upload ${file.name}`);
            return res.json();
          })
          .then(data => {
            if (data.data && data.data.filePath && data.data.thumbnailPath) {
              galleryImages.push({
                image: data.data.filePath,
                thumbnail: data.data.thumbnailPath
              });
            } else {
              throw new Error(`Invalid response for ${file.name}`);
            }
          });
      });

      Promise.all(uploadPromises)
        .then(() => {
          renderGallery();
        })
        .catch(err => {
          console.error('Gallery upload error:', err);
          alert(`Some uploads failed: ${err.message}`);
          renderGallery();
        })
        .finally(() => {
          isUploading = false;
          galleryUpload.disabled = false;
          galleryUpload.value = ''; // clear input
        });
    });
  }
});

function renderGallery() {
  const container = document.getElementById('galleryPreview');
  const hiddenInput = document.getElementById('portfolioImagesJSON');
  if (!container || !hiddenInput) return;

  container.innerHTML = '';
  galleryImages.forEach((imgObj, idx) => {
    const item = document.createElement('div');
    item.className = 'gallery-item';
    item.style.position = 'relative';
    item.style.width = '120px';
    item.style.height = '120px';
    item.style.borderRadius = 'var(--border-radius)';
    item.style.overflow = 'hidden';
    item.style.border = '1px solid var(--border-color)';

    const img = document.createElement('img');
    img.src = imgObj.thumbnail;
    img.style.width = '100%';
    img.style.height = '100%';
    img.style.objectFit = 'cover';

    const delBtn = document.createElement('button');
    delBtn.type = 'button';
    delBtn.className = 'gallery-item-delete';
    delBtn.innerHTML = '&times;';
    delBtn.style.position = 'absolute';
    delBtn.style.top = '5px';
    delBtn.style.right = '5px';
    delBtn.style.background = 'rgba(0, 0, 0, 0.6)';
    delBtn.style.color = '#fff';
    delBtn.style.border = 'none';
    delBtn.style.borderRadius = '50%';
    delBtn.style.width = '24px';
    delBtn.style.height = '24px';
    delBtn.style.display = 'flex';
    delBtn.style.alignItems = 'center';
    delBtn.style.justifyContent = 'center';
    delBtn.style.cursor = 'pointer';
    delBtn.style.fontSize = '14px';
    delBtn.addEventListener('click', function () {
      galleryImages.splice(idx, 1);
      renderGallery();
    });

    item.appendChild(img);
    item.appendChild(delBtn);
    container.appendChild(item);
  });

  hiddenInput.value = JSON.stringify(galleryImages);
}

document.getElementById('portfolio-form').addEventListener('submit', async function (e) {
  e.preventDefault();

  if (isUploading) {
    alert('Please wait for image upload to complete before submitting.');
    return;
  }

  let portfolioId = null;
  const pathMatch = window.location.pathname.match(/\/portfolio\/(\d+)\/edit$/);
  if (pathMatch) {
    portfolioId = pathMatch[1];
  }

  const title = document.getElementById('title').value.trim();
  const slug = document.getElementById('slug').value.trim();
  const shortDescription = document.getElementById('shortDescription').value.trim();
  const projectURL = document.getElementById('projectURL').value.trim();
  const githubURL = document.getElementById('githubURL').value.trim();
  const sortOrder = parseInt(document.getElementById('sortOrder').value) || 0;
  const showcaseImage = document.getElementById('showcaseImageURL').value.trim();
  const images = document.getElementById('portfolioImagesJSON').value.trim();

  if (!title || !slug) {
    alert('Please enter a title and slug.');
    return;
  }
  if (!shortDescription) {
    alert('Please enter a short description.');
    return;
  }

  if (projectURL && !isValidUrl(projectURL)) {
    alert('Please enter a valid project URL.');
    return;
  }
  if (githubURL && !isValidUrl(githubURL)) {
    alert('Please enter a valid GitHub URL.');
    return;
  }

  const portfolioData = {
    title: title,
    slug: slug,
    short_description: shortDescription,
    project_url: projectURL || null,
    github_url: githubURL || null,
    showcase_image: showcaseImage || null,
    sort_order: sortOrder,
    images: images || '[]'
  };

  const isEdit = window.location.pathname.includes('/edit');
  const url = isEdit ? `${window.ROOT_PREFIX || ''}/api/portfolio/${portfolioId}` : `${window.ROOT_PREFIX || ''}/api/portfolio`;
  const method = isEdit ? 'PUT' : 'POST';
  const successMessage = isEdit ? 'Portfolio item updated successfully!' : 'Portfolio item created successfully!';

  try {
    const response = await fetch(url, {
      method: method,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(portfolioData)
    });

    if (response.ok) {
      alert(successMessage);
      window.location.href = `${window.ROOT_PREFIX || ''}/portfolio`;
    } else {
      const errorData = await response.json();
      alert(errorData.error || 'An error occurred');
    }
  } catch (error) {
    console.error('Error:', error);
    alert('Network error. Please try again.');
  }
});

function loadPortfolioItem(id) {
  fetch(`${window.ROOT_PREFIX || ''}/api/portfolio/${id}`)
    .then(response => response.json())
    .then(item => {
      document.getElementById('title').value = item.title;
      document.getElementById('slug').value = item.slug || '';
      document.getElementById('shortDescription').value = item.short_description;
      document.getElementById('projectURL').value = item.project_url || '';
      document.getElementById('githubURL').value = item.github_url || '';
      document.getElementById('sortOrder').value = item.sort_order;
      document.getElementById('showcaseImageURL').value = item.showcase_image || '';

      if (item.images) {
        try {
          galleryImages = JSON.parse(item.images);
        } catch (e) {
          console.error('Failed to parse gallery images:', e);
          galleryImages = [];
        }
        renderGallery();
      }

      if (item.showcase_image) {
        const imagePreview = document.getElementById('imagePreview');
        const previewImg = document.getElementById('previewImg');
        previewImg.src = item.showcase_image;
        imagePreview.classList.remove('hidden');
      }

      document.title = `Edit Portfolio Item: ${item.title}`;
      const heading = document.querySelector('h2');
      if (heading) {
        heading.textContent = `Edit Portfolio Item: ${item.title}`;
      }
    })
    .catch(error => {
      console.error('Error loading portfolio item:', error);
      alert('Error loading portfolio item');
    });
}

function slugify(text) {
  return text
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

function isValidUrl(string) {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
}