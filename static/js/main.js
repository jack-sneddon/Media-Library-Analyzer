document.addEventListener('DOMContentLoaded', function() {
    console.log('Initializing Media Library Analyzer...');

    // Initialize filters first
    const currentYear = new Date().getFullYear();
    document.getElementById('year-end').max = currentYear;
    document.getElementById('year-start').max = currentYear;

    // Add event listeners for all filter controls
    const filterControls = [
        'decade-filter',
        'status-filter',
        'month-filter',
        'file-threshold',
        'year-start',
        'year-end'
    ];

    filterControls.forEach(id => {
        const element = document.getElementById(id);
        if (element) {
            const eventType = element.tagName === 'INPUT' ? 'input' : 'change';
            element.addEventListener(eventType, () => {
                console.log(`Filter changed: ${id}`);
                applyFilters();
            });
        }
    });

    // Initial filter status
    updateActiveFilters(getActiveFilters());
    console.log('Initialization complete');
});

function getActiveFilters() {
    console.log('Getting active filters...');
    const filters = {
        decade: document.getElementById('decade-filter')?.value || 'all',
        status: document.getElementById('status-filter')?.value || 'all',
        month: document.getElementById('month-filter')?.value || 'all',
        fileThreshold: document.getElementById('file-threshold')?.value || 'all',
        yearStart: null,
        yearEnd: null
    };

    const yearStart = document.getElementById('year-start')?.value;
    const yearEnd = document.getElementById('year-end')?.value;

    if (yearStart) filters.yearStart = parseInt(yearStart);
    if (yearEnd) filters.yearEnd = parseInt(yearEnd);

    console.log('Current filters:', filters);
    return filters;
}

function hideElement(element) {
    element.classList.add('hidden');
}

function showElement(element) {
    element.classList.remove('hidden');
}

function setElementVisibility(element, visible) {
    if (visible) {
        showElement(element);
    } else {
        hideElement(element);
    }
}

function applyFilters() {
    const filters = getActiveFilters();
    let visibleCount = 0;

    const years = document.querySelectorAll('.year');
    years.forEach(year => {
        const yearNum = parseInt(year.getAttribute('data-year'));
        let yearVisible = false;

        if (checkYearFilters(yearNum, filters)) {
            const months = year.querySelectorAll('.month-card');
            months.forEach(month => {
                const monthVisible = checkMonthFilters(month, filters);
                setElementVisibility(month, monthVisible);
                if (monthVisible) {
                    yearVisible = true;
                    visibleCount++;
                }
            });
        }

        setElementVisibility(year, yearVisible);
    });

    console.log(`Applied filters, ${visibleCount} items visible`);
    updateActiveFilters(filters);
}

function resetFilters(triggerUpdate = true) {
    // Reset filter values
    document.getElementById('decade-filter').value = 'all';
    document.getElementById('status-filter').value = 'all';
    document.getElementById('month-filter').value = 'all';
    document.getElementById('file-threshold').value = 'all';
    document.getElementById('year-start').value = '';
    document.getElementById('year-end').value = '';

    // Show all elements
    document.querySelectorAll('.year').forEach(showElement);
    document.querySelectorAll('.month-card').forEach(showElement);

    if (triggerUpdate) {
        applyFilters();
    }
}

function checkYearFilters(year, filters) {
    if (filters.decade !== 'all') {
        const decade = Math.floor(year / 10) * 10;
        if (decade !== parseInt(filters.decade)) {
            return false;
        }
    }

    if (filters.yearStart && year < filters.yearStart) return false;
    if (filters.yearEnd && year > filters.yearEnd) return false;

    return true;
}

function checkMonthFilters(monthCard, filters) {
    const status = monthCard.getAttribute('data-status');
    const fileCount = parseInt(monthCard.querySelector('.file-count').textContent.match(/\d+/)[0]);
    const monthName = monthCard.querySelector('.month-name').textContent.trim();

    // Status check
    if (filters.status !== 'all' && status !== filters.status) {
        return false;
    }

    // Month check
    if (filters.month !== 'all' && !monthName.toLowerCase().includes(filters.month.toLowerCase())) {
        return false;
    }

    // Threshold check
    if (filters.fileThreshold !== 'all') {
        switch (filters.fileThreshold) {
            case '0':
                if (fileCount !== 0) return false;
                break;
            case '1-29':
                if (fileCount < 1 || fileCount > 29) return false;
                break;
            case '30+':
                if (fileCount < 30) return false;
                break;
        }
    }

    return true;
}

function applyPreset(presetName) {
    console.log(`Applying preset: ${presetName}`);
    const currentYear = new Date().getFullYear();

    resetFilters(false);

    switch(presetName) {
        case 'recent-missing':
            document.getElementById('year-start').value = currentYear - 1;
            document.getElementById('year-end').value = currentYear;
            document.getElementById('status-filter').value = 'missing';
            break;
        case 'all-missing':
            document.getElementById('status-filter').value = 'missing';
            break;
        case 'incomplete':
            document.getElementById('file-threshold').value = '1-29';
            break;
    }

    applyFilters();
}


function updateActiveFilters(filters) {
    if (!filters) return;

    const activeFilters = [];
    
    if (filters.decade !== 'all') activeFilters.push(`Decade: ${filters.decade}s`);
    if (filters.status !== 'all') activeFilters.push(`Status: ${filters.status}`);
    if (filters.month !== 'all') activeFilters.push(`Month: ${filters.month}`);
    if (filters.fileThreshold !== 'all') activeFilters.push(`Files: ${filters.fileThreshold}`);
    if (filters.yearStart) activeFilters.push(`From: ${filters.yearStart}`);
    if (filters.yearEnd) activeFilters.push(`To: ${filters.yearEnd}`);

    const filterDisplay = document.getElementById('active-filters');
    if (filterDisplay) {
        filterDisplay.textContent = activeFilters.length > 0 
            ? 'Active Filters: ' + activeFilters.join(' | ')
            : 'No active filters';
    }
}