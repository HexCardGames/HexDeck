import { addMessages, register, init, getLocaleFromNavigator } from 'svelte-i18n';

import en from './en.json';
import de from './de.json';

addMessages('en', en);
addMessages('de', de);

register('en', () => import('./en.json'));
register('de', () => import('./de.json'));

const initialLocale = getLocaleFromNavigator();
console.log('Initial locale:', initialLocale);

init({
  fallbackLocale: 'en',
  initialLocale: initialLocale,
});