import React from 'react';
import OriginalTOC from '@theme-original/TOC';
import EditThisPage from '@theme/EditThisPage';

export default function TOC({ toc, editUrl, ...props }) : JSX.Element {
  return (
    <>
      <h2>Contents</h2>
      <OriginalTOC toc={toc} {...props} />
      <EditThisPage editUrl={editUrl} />
    </>
  );
}
