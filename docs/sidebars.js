const sidebars = {
  docsSidebar: [
    {
      type: 'doc',
      id: 'introduction',
      label: 'Introduction',
    },
    {
      type: 'doc',
      id: 'install',
      label: 'Install',
    },
    {
      type: 'category',
      label: 'Command Line Usage',
      items: [
        'cmd/root',
        'cmd/tag',
        'cmd/bump',
        'cmd/changelog',
        'cmd/release',
      ],
    },
    {
      type: 'category',
      label: 'Configuration',
      items: [
        'config/about',
        'config/basics',
        'config/bumping',
        'config/changelog',
      ],
    },
  ],
};

module.exports = sidebars;
