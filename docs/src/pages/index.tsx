import Head from '@docusaurus/Head';
import { Redirect } from '@docusaurus/router';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import React from 'react';

export default function HomePage(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <>
      <Head>
        <meta title={siteConfig.title} />
        <meta property="og:description" content={siteConfig.tagline} />
        <meta property="description" content={siteConfig.tagline} />
        <link rel="canonical" href="https://github.com/gembaadvantage/uplift" />
      </Head>
      <Redirect to="/docs/introduction" />
    </>
  );
}
