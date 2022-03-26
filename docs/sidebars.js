const sidebars = {
  docsSidebar: [
    {
      type: "category",
      label: "Home",
      collapsible: false,
      items: [
        "home/introduction",
        {
          type: "link",
          label: "Installation",
          href: "#installation",
        },
        { type: "link", label: "Quick Start", href: "#quick-start" },
      ],
    },
    {
      type: "category",
      label: "Command Line Usage",
      items: [
        "cmd/root",
        "cmd/tag",
        "cmd/bump",
        "cmd/changelog",
        "cmd/release",
      ],
    },
    {
      type: "category",
      label: "Configuration",
      items: [
        "config/about",
        "config/basics",
        "config/bumping",
        "config/changelog",
      ],
    },
  ],
};

module.exports = sidebars;
