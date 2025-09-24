function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
    tab.classList.add('hidden');
    });
    document.querySelectorAll('.tab-button').forEach(btn => {
    btn.classList.remove('active', 'bg-white', 'dark:bg-gray-600', 'text-blue-600', 'dark:text-blue-400', 'shadow-sm');
    btn.classList.add('text-gray-500', 'dark:text-gray-400');
    });

    // Show selected tab
    document.getElementById(tabName + '-tab').classList.remove('hidden');
    event.target.classList.add('active', 'bg-white', 'dark:bg-gray-600', 'text-blue-600', 'dark:text-blue-400', 'shadow-sm');
    event.target.classList.remove('text-gray-500', 'dark:text-gray-400');
}