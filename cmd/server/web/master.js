// Master Data Management

// Navigation
document.querySelectorAll('.nav-item').forEach(item => {
  item.addEventListener('click', (e) => {
    e.preventDefault();
    const page = item.getAttribute('data-page');
    switchPage(page);
  });
});

function switchPage(pageName) {
  // Update nav active state
  document.querySelectorAll('.nav-item').forEach(item => {
    item.classList.remove('active');
    if (item.getAttribute('data-page') === pageName) {
      item.classList.add('active');
    }
  });

  // Update page content
  document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
  document.getElementById('page-' + pageName).classList.add('active');

  // Update header and button
  const pageTitle = document.getElementById('pageTitle');
  const btnAdd = document.getElementById('btnAdd');

  if (pageName === 'logs') {
    pageTitle.textContent = 'Activity Logs';
    btnAdd.innerHTML = '<span>➕</span><span>Add Log</span>';
    btnAdd.onclick = () => { openModal('Add Log'); };
    btnAdd.style.display = 'flex';
  } else if (pageName.startsWith('master-')) {
    const masterType = pageName.replace('master-', '').replace('-', ' ');
    const titleCase = masterType.split(' ').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ');
    pageTitle.textContent = 'Master ' + titleCase;
    btnAdd.style.display = 'none';
    loadMasterPage(pageName);
  }
}

// Load master page content
function loadMasterPage(pageName) {
  const page = document.getElementById('page-' + pageName);
  const masterType = pageName.replace('master-', '');
  
  page.innerHTML = `
    <div class="master-table-wrap">
      <div class="master-actions">
        <h3>${getMasterTitle(masterType)}</h3>
        <button class="btn btn-primary" onclick="openMasterModal('${masterType}', 'add')">
          ➕ Add ${getMasterSingular(masterType)}
        </button>
      </div>
      <table>
        <thead>
          <tr>
            ${getMasterHeaders(masterType)}
            <th style="text-align:center;">Action</th>
          </tr>
        </thead>
        <tbody id="master-${masterType}-body"></tbody>
      </table>
    </div>
  `;
  
  loadMasterData(masterType);
}

function getMasterTitle(type) {
  const titles = {
    'job-titles': 'Job Titles',
    'pics': 'Person In Charge (PIC)',
    'statuses': 'Status Options',
    'categories': 'Category Options'
  };
  return titles[type] || type;
}

function getMasterSingular(type) {
  const singulars = {
    'job-titles': 'Job Title',
    'pics': 'PIC',
    'statuses': 'Status',
    'categories': 'Category'
  };
  return singulars[type] || type;
}

function getMasterHeaders(type) {
  const headers = {
    'job-titles': '<th style="text-align:center;">ID</th><th>Name</th>',
    'pics': '<th style="text-align:center;">ID</th><th>Name</th><th>Email</th>',
    'statuses': '<th style="text-align:center;">ID</th><th>Name</th><th>Color</th>',
    'categories': '<th style="text-align:center;">ID</th><th>Name</th><th>Description</th>'
  };
  return headers[type] || '<th>ID</th><th>Name</th>';
}

// Load master data from API
async function loadMasterData(type) {
  try {
    const res = await fetch(`/api/master/${type}`);
    const data = await res.json();
    renderMasterData(type, data);
  } catch (error) {
    console.error('Failed to load master data:', error);
  }
}

function renderMasterData(type, data) {
  const tbody = document.getElementById(`master-${type}-body`);
  tbody.innerHTML = '';

  if (!data || data.length === 0) {
    tbody.innerHTML = '<tr><td colspan="10" style="text-align:center;padding:40px;color:var(--gray-500);">No data available</td></tr>';
    return;
  }

  data.forEach(item => {
    const tr = document.createElement('tr');
    tr.innerHTML = getMasterRowHTML(type, item);
    tbody.appendChild(tr);
  });
}

function getMasterRowHTML(type, item) {
  let cells = `<td style="text-align:center;">${item.id}</td>`;
  
  if (type === 'job-titles') {
    cells += `<td>${escapeHtml(item.name)}</td>`;
  } else if (type === 'pics') {
    cells += `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.email || '')}</td>`;
  } else if (type === 'statuses') {
    cells += `<td>${escapeHtml(item.name)}</td><td><span style="color:${item.color}">${item.color || '-'}</span></td>`;
  } else if (type === 'categories') {
    cells += `<td>${escapeHtml(item.name)}</td><td>${escapeHtml(item.description || '')}</td>`;
  }

  cells += `
    <td style="text-align:center;">
      <button class="btn btn-sm" onclick='editMasterItem("${type}", ${JSON.stringify(item)})'>✏️ Edit</button>
      <button class="btn btn-sm btn-danger" onclick="deleteMasterItem('${type}', ${item.id})">🗑️ Delete</button>
    </td>
  `;
  
  return cells;
}

// Master modal
let currentMasterType = '';
let currentMasterItem = null;

function openMasterModal(type, mode, item = null) {
  currentMasterType = type;
  currentMasterItem = item;
  
  const singular = getMasterSingular(type);
  const title = mode === 'add' ? `Add ${singular}` : `Edit ${singular}`;
  
  const formHTML = getMasterFormHTML(type, item);
  
  const modal = document.getElementById('masterModal');
  const content = document.getElementById('masterModalContent');
  content.innerHTML = `
    <h2>${title}</h2>
    <form id="masterForm">
      ${formHTML}
      <div class="modal-actions">
        <button type="button" class="btn btn-ghost" onclick="closeMasterModal()">Cancel</button>
        <button type="submit" class="btn btn-primary">Save</button>
      </div>
    </form>
  `;
  modal.classList.remove('hidden');
  document.getElementById('masterForm').addEventListener('submit', saveMasterItem);
}

function getMasterFormHTML(type, item) {
  if (type === 'job-titles') {
    return `<label>Name <input type="text" name="name" value="${item?.name || ''}" required></label>`;
  } else if (type === 'pics') {
    return `
      <label>Name <input type="text" name="name" value="${item?.name || ''}" required></label>
      <label>Email <input type="email" name="email" value="${item?.email || ''}"></label>
    `;
  } else if (type === 'statuses') {
    return `
      <label>Name <input type="text" name="name" value="${item?.name || ''}" required></label>
      <label>Color <input type="color" name="color" value="${item?.color || '#000000'}"></label>
    `;
  } else if (type === 'categories') {
    return `
      <label>Name <input type="text" name="name" value="${item?.name || ''}" required></label>
      <label>Description <textarea name="description">${item?.description || ''}</textarea></label>
    `;
  }
}

function closeMasterModal() {
  document.getElementById('masterModal').classList.add('hidden');
}

async function saveMasterItem(e) {
  e.preventDefault();
  const form = e.target;
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());
  
  try {
    const url = currentMasterItem 
      ? `/api/master/${currentMasterType}/${currentMasterItem.id}`
      : `/api/master/${currentMasterType}`;
    const method = currentMasterItem ? 'PUT' : 'POST';
    
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    
    if (!res.ok) throw new Error('Failed to save');
    
    closeMasterModal();
    loadMasterData(currentMasterType);
  } catch (error) {
    alert('Failed to save: ' + error.message);
  }
}

function editMasterItem(type, item) {
  openMasterModal(type, 'edit', item);
}

async function deleteMasterItem(type, id) {
  if (!confirm('Delete this item?')) return;
  
  try {
    const res = await fetch(`/api/master/${type}/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error('Failed to delete');
    loadMasterData(type);
  } catch (error) {
    alert('Failed to delete: ' + error.message);
  }
}

// Refresh dropdowns in main form
async function loadDropdownOptions() {
  try {
    const [jobTitles, pics, statuses, categories] = await Promise.all([
      fetch('/api/master/job-titles').then(r => r.json()),
      fetch('/api/master/pics').then(r => r.json()),
      fetch('/api/master/statuses').then(r => r.json()),
      fetch('/api/master/categories').then(r => r.json())
    ]);
    
    // Update filter dropdowns
    updateDropdown('fPic', pics, 'All PIC');
    updateDropdown('fStatus', statuses, 'All Status');
    updateDropdown('fCategory', categories, 'All Category');
    
    // Update form dropdowns
    updateDropdown('fldJobTitle', jobTitles);
    updateDropdown('fldPic', pics);
    updateDropdown('fldStatus', statuses);
    updateDropdown('fldCategory', categories);
  } catch (error) {
    console.error('Failed to load dropdown options:', error);
  }
}

function updateDropdown(id, items, firstOption = null) {
  const select = document.getElementById(id);
  if (!select) return;
  
  let html = firstOption ? `<option value="">${firstOption}</option>` : '';
  items.forEach(item => {
    html += `<option value="${escapeHtml(item.name)}">${escapeHtml(item.name)}</option>`;
  });
  select.innerHTML = html;
}

// Load dropdowns on page load
loadDropdownOptions();
