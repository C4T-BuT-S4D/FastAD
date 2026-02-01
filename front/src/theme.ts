import '@mui/material/styles';
import { createTheme } from '@mui/material/styles';

declare module '@mui/material/styles' {
  // eslint-disable-next-line no-unused-vars
  interface Theme {
    name: string;
  }

  // eslint-disable-next-line no-unused-vars
  interface ThemeOptions {
    name: string;
  }

  interface Palette {
    statusUp: Palette['primary'];
    statusCorrupt: Palette['primary'];
    statusMumble: Palette['primary'];
    statusDown: Palette['primary'];
    statusCheckFailed: Palette['primary'];
    statusUnspecified: Palette['primary'];
  }

  interface PaletteOptions {
    statusUp: PaletteOptions['primary'];
    statusCorrupt: PaletteOptions['primary'];
    statusMumble: PaletteOptions['primary'];
    statusDown: PaletteOptions['primary'];
    statusCheckFailed: PaletteOptions['primary'];
    statusUnspecified: PaletteOptions['primary'];
  }
}

let theme = createTheme({ name: 'ForcAD' });

theme = createTheme(theme, {
  palette: {
    statusUp: theme.palette.augmentColor({
      color: {
        main: '#7dfc74',
      },
      name: 'status_up',
    }),
    statusCorrupt: theme.palette.augmentColor({
      color: {
        main: '#5191ff',
      },
      name: 'status_corrupt',
    }),
    statusMumble: theme.palette.augmentColor({
      color: {
        main: '#ff9114',
      },
      name: 'status_mumble',
    }),
    statusDown: {
      ...theme.palette.augmentColor({
        color: {
          main: '#ff5b5b',
        },
        name: 'status_down',
      }),
      contrastText: '#000000',
    },
    statusCheckFailed: {
      ...theme.palette.augmentColor({
        color: {
          main: '#ffff00',
        },
        name: 'status_check_failed',
      }),
      contrastText: '#000000',
    },
    statusUnspecified: {
      ...theme.palette.augmentColor({
        color: {
          main: '#ffffff',
        },
        name: 'status_unspecified',
      }),
    },
  },
});

export default theme;
