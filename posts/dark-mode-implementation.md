---
title: "Implementing Dark Mode with Tailwind CSS"
date: "2025-01-18"
tags: ["css", "tailwind", "frontend", "design"]
---

# Implementing Dark Mode with Tailwind CSS

Dark mode has become an essential feature for modern web applications. Here's how I implemented it using Tailwind CSS.

## Why Dark Mode?

- **Reduced eye strain** in low-light environments
- **Better battery life** on OLED displays
- **User preference** for different times of day
- **Modern aesthetic** appeal

## Implementation Steps

### 1. Configure Tailwind for Dark Mode

```javascript
// tailwind.config.js
module.exports = {
  darkMode: 'class', // or 'media'
  // ... other config
}
```

### 2. HTML Structure

Start with a class on the HTML element:

```html
<html class="light">
```

### 3. CSS Classes

Use Tailwind's dark: prefix:

```html
<div class="bg-white dark:bg-gray-900">
  <h1 class="text-gray-900 dark:text-white">
    Dark Mode Example
  </h1>
</div>
```

### 4. JavaScript Toggle

```javascript
const html = document.documentElement;
const themeToggle = document.getElementById('theme-toggle');

// Check saved preference
const currentTheme = localStorage.getItem('theme') || 'light';
html.classList.toggle('dark', currentTheme === 'dark');

// Toggle handler
themeToggle.addEventListener('click', () => {
  const isDark = html.classList.toggle('dark');
  localStorage.setItem('theme', isDark ? 'dark' : 'light');
});
```

## Best Practices

1. **Consistent color scheme** across all elements
2. **Smooth transitions** for theme switching
3. **Persistent user preference** with localStorage
4. **System preference** detection as fallback
5. **Accessible colors** with proper contrast ratios

## Common Pitfalls

- **Forgetting borders**: Remember `dark:border-gray-700`
- **Icons and SVGs**: Use `dark:text-white` classes
- **Third-party components**: May need custom CSS
- **Print styles**: Ensure readability in both themes

## Conclusion

Implementing dark mode with Tailwind CSS is straightforward and provides a great user experience. The key is consistency and thorough testing across all UI components.

---

Happy styling! 🎨