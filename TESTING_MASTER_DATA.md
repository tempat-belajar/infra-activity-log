# Master Data Testing Guide

## Overview
This update adds master data management functionality with:
- Sidebar navigation
- CRUD operations for Job Titles, PICs, Statuses, and Categories
- Dynamic form dropdowns from database
- Purple-themed UI matching office standard

## Testing Steps

### 1. Build and Deploy
```bash
# The CI/CD pipeline will automatically:
# - Run migrations (including 0003_master_tables.sql)
# - Build the application
# - Deploy via docker compose
```

### 2. Verify Database Migration
```bash
# SSH to server
docker exec -it infra-activity-log-postgres-1 psql -U akr -d activitylog

# Check if master tables exist
\dt master_*

# Verify default data
SELECT * FROM master_job_titles;
SELECT * FROM master_pics;
SELECT * FROM master_statuses;
SELECT * FROM master_categories;
```

Expected results:
- 3 job titles (DBA, DCO, NETWORK)
- 5 PICs (Zaqi, Dwi, Irwan, Kristopelly, Benny)
- 4 statuses (Open, Process, Hold, Done)
- 3 categories (Change, Daily, Incident)

### 3. Test Frontend UI

#### Access the Application
Open browser: http://[SERVER_IP]:8084

#### Test Sidebar Navigation
1. Verify sidebar appears on the left with purple background
2. Click "Activity Logs" - should show activity logs page
3. Click "Job Titles" - should show master job titles page
4. Click "PICs" - should show master PICs page
5. Click "Statuses" - should show master statuses page
6. Click "Categories" - should show master categories page

#### Test Master Job Titles
1. Navigate to Master > Job Titles
2. Verify list shows: DBA, DCO, NETWORK
3. Click "➕ Add Job Title"
4. Enter name: "DEVOPS"
5. Click Save
6. Verify "DEVOPS" appears in list
7. Click "✏️ Edit" on DEVOPS
8. Change name to "DevOps Engineer"
9. Click Save
10. Verify name updated
11. Click "🗑️ Delete" on DevOps Engineer
12. Confirm deletion
13. Verify item removed from list

#### Test Master PICs
1. Navigate to Master > PICs
2. Verify list shows all 5 default PICs
3. Click "➕ Add PIC"
4. Enter name: "Test User"
5. Enter email: "test@example.com"
6. Click Save
7. Verify new PIC appears with email
8. Test Edit functionality
9. Test Delete functionality

#### Test Master Statuses
1. Navigate to Master > Statuses
2. Verify list shows 4 default statuses
3. Click "➕ Add Status"
4. Enter name: "Pending"
5. Select color: #FFA500 (orange)
6. Click Save
7. Verify color displays correctly
8. Test Edit and Delete

#### Test Master Categories
1. Navigate to Master > Categories
2. Verify list shows 3 default categories
3. Click "➕ Add Category"
4. Enter name: "Maintenance"
5. Enter description: "Regular maintenance tasks"
6. Click Save
7. Verify description appears
8. Test Edit and Delete

### 4. Test Dynamic Dropdowns in Activity Log Form

1. Navigate to Activity Logs page
2. Click "➕ Add Log" button
3. Verify dropdowns are populated:
   - Job Title dropdown shows data from master_job_titles
   - PIC dropdown shows data from master_pics
   - Status dropdown shows data from master_statuses
   - Category dropdown shows data from master_categories
4. Create a new log entry using dropdown values
5. Verify log is created successfully

### 5. Test Filter Dropdowns
1. On Activity Logs page, verify filter dropdowns also populated from master data
2. Test filtering by PIC, Status, and Category
3. Verify filters work correctly

## API Endpoints to Test

### Master Job Titles
- GET    /api/master/job-titles       - List all
- POST   /api/master/job-titles       - Create
- PUT    /api/master/job-titles/:id   - Update
- DELETE /api/master/job-titles/:id   - Delete (soft delete)

### Master PICs
- GET    /api/master/pics       - List all
- POST   /api/master/pics       - Create
- PUT    /api/master/pics/:id   - Update
- DELETE /api/master/pics/:id   - Delete (soft delete)

### Master Statuses
- GET    /api/master/statuses       - List all
- POST   /api/master/statuses       - Create
- PUT    /api/master/statuses/:id   - Update
- DELETE /api/master/statuses/:id   - Delete (soft delete)

### Master Categories
- GET    /api/master/categories       - List all
- POST   /api/master/categories       - Create
- PUT    /api/master/categories/:id   - Update
- DELETE /api/master/categories/:id   - Delete (soft delete)

## Expected Behavior

### Soft Delete
- Deleted items set `is_active = false`
- Items with `is_active = false` don't appear in lists
- Data preserved in database for audit purposes

### Validation
- Name fields are required
- Duplicate names should be prevented (database UNIQUE constraint)
- Form validation prevents empty submissions

## Troubleshooting

### Migration Issues
```bash
# Check migration status
docker logs infra-activity-log-postgres-1 | grep migration

# Manually run migration if needed
docker exec -i infra-activity-log-postgres-1 psql -U akr -d activitylog < migrations/0003_master_tables.sql
```

### API Errors
```bash
# Check application logs
docker logs infra-activity-log-app-1 | tail -50

# Check if routes are registered
docker logs infra-activity-log-app-1 | grep "master"
```

### Frontend Issues
- Open browser DevTools (F12)
- Check Console for JavaScript errors
- Check Network tab for failed API requests
- Verify all JS files loaded: app.js and master.js

## Success Criteria
✅ All 4 master data pages accessible via sidebar
✅ CRUD operations work for all master types
✅ Form dropdowns populated from master data
✅ Filter dropdowns populated from master data
✅ No console errors in browser
✅ No errors in application logs
✅ Database contains master tables with data
