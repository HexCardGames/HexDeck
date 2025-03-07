import { writable } from "svelte/store";

const getSystemTheme = () => (window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light");

const storedTheme = localStorage.getItem("theme") as Theme | null;
const initialTheme: Theme = storedTheme === "dark" || storedTheme === "light" ? storedTheme : "system";

type Theme = "dark" | "light" | "system";
export const theme = writable<Theme>(initialTheme);

const applyTheme = (value: Theme) => {
    const resolvedTheme = value === "system" ? getSystemTheme() : value;

    document.documentElement.classList.toggle("dark", resolvedTheme === "dark");
    document.documentElement.setAttribute("data-theme", resolvedTheme);
    document.body.classList.toggle("dark-theme", resolvedTheme === "dark");
    document.body.classList.toggle("light-theme", resolvedTheme === "light");

    localStorage.setItem("theme", value);
};

theme.subscribe(applyTheme);

// Watch for system theme changes when "system" mode is enabled
const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
mediaQuery.addEventListener("change", () => {
    theme.update((current) => {
        if (current === "system") applyTheme("system");
        return current;
    });
});

export const setTheme = (value: Theme) => {
    theme.set(value);
};

export const toggleTheme = () => {
    theme.update((current) => {
        if (current === "dark") return "light";
        if (current === "light") return "system";
        return "dark";
    });
};
