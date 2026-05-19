// Global variables
let isUploading = false;

document.addEventListener('DOMContentLoaded', function () {
  // Check if editing
  let portfolioId = null;
  const pathMatch = window.location.pathname.match(/^\/admin\/portfolio\/(\d+)\/edit$/);
  if (pathMatch) {
    portfolioId = pathMatch[1];
    loadPortfolioItem(portfolioId);
  }

  // Set up image upload handler
  const imageUpload = document.getElementById('showcaseImageUpload');
  const imageUrlInput = document.getElementById('showcaseImageURL');
  const imagePreview = document.getElementById('imagePreview');
  const previewImg = document.getElementById('previewImg');

  imageUpload.addEventListener('change', function (e) {
    const file = e.target.files[0];
    if (file) {
      // Validate file type
      if (!file.type.startsWith('image/')) {
        alert('Please select an image file');
        return;
      }

      // Validate file size (5MB limit)
      if (file.size > 5 * 1024 * 1024) {
        alert('File size must be less than 5MB');
        return;
      }

      // Set uploading state
      isUploading = true;
      imageUpload.disabled = true;

      // Upload the file
      const formData = new FormData();
      formData.append('image', file);

      fetch('/api/upload/image', {
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
});

document.getElementById('portfolio-form').addEventListener('submit', async function (e) {
  e.preventDefault();

  // Prevent submission if upload is in progress
  if (isUploading) {
    alert('Please wait for image upload to complete before submitting.');
    return;
  }

  // Check if editing
  let portfolioId = null;
  const pathMatch = window.location.pathname.match(/^\/admin\/portfolio\/(\d+)\/edit$/);
  if (pathMatch) {
    portfolioId = pathMatch[1];
  }

  const title = document.getElementById('title').value.trim();
  const shortDescription = document.getElementById('shortDescription').value.trim();
  const projectURL = document.getElementById('projectURL').value.trim();
  const githubURL = document.getElementById('githubURL').value.trim();
  const sortOrder = parseInt(document.getElementById('sortOrder').value) || 0;
  const showcaseImage = document.getElementById('showcaseImageURL').value.trim();

  // Basic validation
  if (!title) {
    alert('Please enter a title.');
    return;
  }
  if (!shortDescription) {
    alert('Please enter a short description.');
    return;
  }

  // URL validation
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
    short_description: shortDescription,
    project_url: projectURL || null,
    github_url: githubURL || null,
    showcase_image: showcaseImage || null,
    sort_order: sortOrder
  };

  const isEdit = window.location.pathname.includes('/edit');
  const url = isEdit ? `/api/portfolio/${portfolioId}` : '/api/portfolio';
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
      window.location.href = '/admin/portfolio';
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
  fetch(`/api/portfolio/${id}`)
    .then(response => response.json())
    .then(item => {
      document.getElementById('title').value = item.title;
      document.getElementById('shortDescription').value = item.short_description;
      document.getElementById('projectURL').value = item.project_url || '';
      document.getElementById('githubURL').value = item.github_url || '';
      document.getElementById('sortOrder').value = item.sort_order;
      document.getElementById('showcaseImageURL').value = item.showcase_image || '';

      // Show image preview if exists
      if (item.showcase_image) {
        const imagePreview = document.getElementById('imagePreview');
        const previewImg = document.getElementById('previewImg');
        previewImg.src = item.showcase_image;
        imagePreview.classList.remove('hidden');
      }

      // Update page title
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

function isValidUrl(string) {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
}