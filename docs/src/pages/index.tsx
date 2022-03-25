import React from 'react';
import Head from '@docusaurus/Head';
import { Redirect } from '@docusaurus/router';

export default function HomePage(): JSX.Element {
  return (
    <>
      <Head>
        <meta title="Uplift"/>
        <meta
          property="og:description"
          content="Semantic versioning the easy way. Powered by Conventional Commits. Built for use with CI"
        />
        <meta
          property="description"
          content="Semantic versioning the easy way. Powered by Conventional Commits. Built for use with CI"
        />
        <link rel="canonical" href="https://github.com/gembaadvantage/uplift" />
      </Head>
      <Redirect to="/docs/home/introduction" />
    </>
  );
}
