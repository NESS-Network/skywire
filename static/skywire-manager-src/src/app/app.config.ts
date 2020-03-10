export const AppConfig = {
  /**
   * How many elements the short lists can have. If the list has more elements, a
   * link for opening the full list is shown.
   */
  maxShortListElements: 5,
  /**
   * How many elements the full lists can have per page.
   */
  maxFullListElements: 40,

  /**
   * Available languages.
   */
  languages: [
    {
      code: 'en',
      name: 'English',
      iconName: 'en.png',
    },
    {
      code: 'es',
      name: 'Español',
      iconName: 'es.png',
    },
  ],
  /**
   * Default language.
   */
  defaultLanguage: 'en',

  /**
   * Sizes of the modal windows.
   */
  smallModalWidth: '480px',
  mediumModalWidth: '640px',
  largeModalWidth: '900px',
};