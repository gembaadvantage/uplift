// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const { tailwindPlugin } = require('./src/plugins');

const config = {
  title: 'Uplift',
  tagline:
    'Semantic versioning the easy way. Powered by Conventional Commits. Built for use with CI',
  url: 'https://upliftci.dev',
  baseUrl: '/',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'favicon.ico',
  organizationName: 'Gemba Advantage',
  projectName: 'uplift',
  clientModules: [require.resolve('./src/css/tailwind.css')],

  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          id: 'default',
          sidebarPath: require.resolve('./sidebars.js'),
          sidebarCollapsible: true,
          breadcrumbs: false,
          editUrl: 'https://github.com/gembaadvantage/uplift/edit/main/docs',
        },
        blog: false,
      },
    ],
  ],

  themeConfig: {
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: true,
    },
    navbar: {
      hideOnScroll: false,
      title: 'Uplift',
      logo: {
        alt: 'My Site Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          href: 'https://github.com/gembaadvantage/uplift',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    hideableSidebar: false,
    footer: {
      copyright: `Copyright © ${new Date().getFullYear()} My Project, Inc. Built with Docusaurus.`,
    },
    prism: {
      theme: require('prism-react-renderer/themes/vsDark'),
    },
  },
  plugins: [tailwindPlugin],
};

module.exports = config;
