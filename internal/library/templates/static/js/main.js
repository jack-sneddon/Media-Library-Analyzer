document.addEventListener('DOMContentLoaded', function() {
    // Initialize filters
    initializeFilters();
    
    // Add event listeners
    document.getElementById('decade-filter').addEventListener('change', applyFilters);
    document.getElementById('status-filter').addEventListener('change', applyFilters);
    document.getElementById('year-start').addEventListener('input', applyFilters);
    document.getElementById('year-end').addEventListener('input', applyFilters);
    document.getElementById('min-files').addEventListener('input', applyFilters);
    document.getElementById('month-filter').addEventListener('change', applyFilters);
});

function initializeFilters() {
    const currentYear = new Date().getFullYear();
    document.getElementById('year-end').max = currentYear;
    document.getElementById('year-start').max = currentYear;
    updateActiveFilters();
}

function applyFilters() {
    const filters = getActiveFilters();
    const years = document.querySelectorAll('.year');
    
    years.forEach(year => {
        const yearNum = parseInt(year.getAttribute('data-year'));
        const shouldShowYear = checkYearFilters(yearNum, filters);
        
        if (shouldShowYear) {
            const months = year.querySelectorAll('.month-card');
            let hasVisibleMonths = false;
            
            months.forEach(month => {
                const shouldShowMonth = checkMonthFilters(month, filters);
                month.classList.toggle('hidden', !shouldShowMonth);
                if (shouldShowMonth) hasVisibleMonths = true;
            });
            
            year.classList.toggle('hidden', !hasVisibleMonths);
        } else {
            year.classList.add('hidden');
        }
    });
    
    updateActiveFilters();
}

function getActiveFilters() {
    return {
        decade: document.getElementById('decade-filter').value,
        status: document.getElementById('status-filter').value,
        yearStart: parseInt(document.getElementById('year-start').value) || 1980,
        yearEnd: parseInt(document.getElementById('year-end').value) || new Date().getFullYear(),
        minFiles: parseInt(document.getElementById('min-files').value) || 0,
        month: document.getElementById('month-filter').value
    };
}

function checkYearFilters(year, filters) {
    const decade = Math.floor(year / 10) * 10;
    return (filters.decade === 'all' || decade === parseInt(filters.decade)) &&
           year >= filters.yearStart &&
           year <= filters.yearEnd;
}

function checkMonthFilters(monthCard, filters) {
    const status = monthCard.getAttribute('data-status');
    const fileCount = parseInt(monthCard.querySelector('.file-count').textContent.match(/\d+/)[0]);
    const monthName = monthCard.querySelector('.month-name').textContent;
    
    return (filters.status === 'all' || status === filters.status) &&
           fileCount >= filters.minFiles &&
           (filters.month === 'all' || monthName.toLowerCase().includes(filters.month.toLowerCase()));
}

function resetFilters() {
    document.getElementById('decade-filter').value = 'all';
    document.getElementById('status-filter').value = 'all';
    document.getElementById('year-start').value = '';
    document.getElementById('year-end').value = '';
    document.getElementById('min-files').value = '';
    document.getElementById('month-filter').value = 'all';
    applyFilters();
}

function updateActiveFilters() {
    const filters = getActiveFilters();
    const activeFilters = [];
    
    if (filters.decade !== 'all') activeFilters.push(`Decade: ${filters.decade}s`);
    if (filters.status !== 'all') activeFilters.push(`Status: ${filters.status}`);
    if (filters.minFiles > 0) activeFilters.push(`Min Files: ${filters.minFiles}`);
    if (filters.month !== 'all') activeFilters.push(`Month: ${filters.month}`);
    if (filters.yearStart !== 1980 || filters.yearEnd !== new Date().getFullYear()) {
        activeFilters.push(`Years: ${filters.yearStart}-${filters.yearEnd}`);
    }
    
    const filterDisplay = document.getElementById('active-filters');
    filterDisplay.textContent = activeFilters.length > 0 
        ? 'Active Filters: ' + activeFilters.join(' | ')
        : 'No active filters';
}