/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
import { useColorMode, useThemeConfig } from '@docusaurus/theme-common';
import { MoonIcon, SunIcon } from '@heroicons/react/outline';
import React from 'react';

export default function NavbarColorModeToggle() : JSX.Element {
  const disabled = useThemeConfig().colorMode.disableSwitch;

  if (disabled) {
    return null;
  }

  const { isDarkTheme, setLightTheme, setDarkTheme } = useColorMode();
  const toggleTheme = () => (isDarkTheme ? setLightTheme() : setDarkTheme());

  return (
    <button
      className="fixed bottom-[1.3rem] right-[1.3rem] z-50 h-12 w-12 cursor-pointer rounded-full bg-primary p-2 text-gray-100 shadow-md transition hover:text-white hover:shadow-xl"
      onClick={toggleTheme}
    >
      {isDarkTheme ? (
        <SunIcon className="h-full w-full text-current" />
      ) : (
        <MoonIcon className="h-full w-full text-current" />
      )}
    </button>
  );
}
