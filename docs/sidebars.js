const sidebars = {
  docsSidebar: [
    {
      type: "category",
      label: "Home",
      items: [
        "home/introduction",
        { type: "link", label: "Installation", href: "#installation" },
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
  ],
};

module.exports = sidebars;
