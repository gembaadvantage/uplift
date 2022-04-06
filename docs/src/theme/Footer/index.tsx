/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
import React from 'react';
import { Github } from '@styled-icons/boxicons-logos';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';

export default function Footer() : JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <footer className="bg-background-100 py-6 lg:px-12">
      <div className="mx-auto flex w-full max-w-6xl flex-row px-8 items-start">
        <div className="basis-3/4">
          <p className="text-sm leading-none">
            Uplift is built and maintained by the folk at <Link href={`https://github.com/${siteConfig.organizationName}`}><img className="relative top-1" src={require('@site/static/img/ga_icon_small.png').default} /></Link>
          </p>
          <p className="text-sm leading-none">
            Made with ❤️ using Docusaurus
          </p>
        </div>
        <div className="basis-1/4">
          <div className="flex flex-row-reverse align-middle">
            <Link
              href={`https://github.com/${siteConfig.organizationName}/${siteConfig.projectName}`}
              className="inline-flex text-current"
            >
              <Github className="h-7" />
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
