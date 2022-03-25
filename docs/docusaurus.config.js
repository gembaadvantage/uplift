// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

const config = {
  title: "Uplift",
  tagline:
    "Semantic versioning the easy way. Powered by Conventional Commits. Built for use with CI",
  url: "https://upliftci.dev",
  baseUrl: "/",
  onBrokenLinks: "warn",
  onBrokenMarkdownLinks: "warn",
  favicon: "favicon.ico",
  organizationName: "Gemba Advantage",
  projectName: "uplift",

  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          id: "default",
          sidebarPath: require.resolve("./sidebars.js"),
          sidebarCollapsible: true,
          breadcrumbs: false,
        },
        blog: false,
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        hideOnScroll: true,
        title: "Uplift",
        logo: {
          alt: "My Site Logo",
          src: "img/logo.svg",
        },
        items: [
          {
            href: "https://github.com/gembaadvantage/uplift",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      hideableSidebar: true,
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
