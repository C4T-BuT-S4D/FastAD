import { useTheme } from '@mui/material/styles';
import ForcADLayout from './layouts/ForcAD.tsx';

function App() {
  const theme = useTheme();

  switch (theme.name) {
    case 'ForcAD':
      return <ForcADLayout />;
    default:
      return <ForcADLayout />;
  }
}

export default App;
