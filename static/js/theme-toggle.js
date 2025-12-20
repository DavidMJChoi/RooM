/**
 * Theme Toggle Module
 * Manages dark/light theme switching functionality
 */
class ThemeToggle {
    constructor() {
        this.themeToggle = document.getElementById('theme-toggle');
        this.sunIcon = document.getElementById('sun-icon');
        this.moonIcon = document.getElementById('moon-icon');
        this.html = document.documentElement;
        
        this.init();
    }

    init() {
        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.setupTheme());
        } else {
            this.setupTheme();
        }
    }

    setupTheme() {
        // Check for saved theme preference or default to light mode
        const currentTheme = localStorage.getItem('theme') || 'light';
        this.html.classList.toggle('dark', currentTheme === 'dark');
        this.updateIcons(currentTheme);

        // Add click event listener to theme toggle button
        if (this.themeToggle) {
            this.themeToggle.addEventListener('click', () => this.toggleTheme());
        }

        // Optional: Respect system preference
        if (currentTheme === 'system' || !localStorage.getItem('theme')) {
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            this.html.classList.toggle('dark', prefersDark);
            this.updateIcons(prefersDark ? 'dark' : 'light');
        }

        // Listen for system theme changes
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
            if (!localStorage.getItem('theme') || localStorage.getItem('theme') === 'system') {
                this.html.classList.toggle('dark', e.matches);
                this.updateIcons(e.matches ? 'dark' : 'light');
            }
        });
    }

    toggleTheme() {
        const isDark = this.html.classList.toggle('dark');
        localStorage.setItem('theme', isDark ? 'dark' : 'light');
        this.updateIcons(isDark ? 'dark' : 'light');
    }

    updateIcons(theme) {
        if (!this.sunIcon || !this.moonIcon) return;
        
        if (theme === 'dark') {
            this.sunIcon.classList.add('hidden');
            this.moonIcon.classList.remove('hidden');
        } else {
            this.sunIcon.classList.remove('hidden');
            this.moonIcon.classList.add('hidden');
        }
    }

    // Public method to get current theme
    getCurrentTheme() {
        return this.html.classList.contains('dark') ? 'dark' : 'light';
    }

    // Public method to set theme programmatically
    setTheme(theme) {
        const validThemes = ['light', 'dark', 'system'];
        if (!validThemes.includes(theme)) {
            console.warn(`Invalid theme "${theme}". Valid themes are: ${validThemes.join(', ')}`);
            return;
        }

        if (theme === 'system') {
            localStorage.removeItem('theme');
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            this.html.classList.toggle('dark', prefersDark);
            this.updateIcons(prefersDark ? 'dark' : 'light');
        } else {
            localStorage.setItem('theme', theme);
            this.html.classList.toggle('dark', theme === 'dark');
            this.updateIcons(theme);
        }
    }
}

// Initialize theme toggle when the module loads
const themeToggle = new ThemeToggle();

// Export for external use if needed
window.ThemeToggle = ThemeToggle;
window.themeToggle = themeToggle;