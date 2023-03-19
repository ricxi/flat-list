import { Routes, Route } from 'react-router-dom';
import UserActivationPage from './pages/UserActivationPage';
import ActivatedSuccessPage from './pages/ActivatedSuccessPage';
import ActivatedFailedPage from './pages/ActivatedFailedPage';

function App() {
  return (
    <main>
      <Routes>
        <Route path='/activate' element={<UserActivationPage />} />
        <Route path='/activated/success' element={<ActivatedSuccessPage />} />
        <Route path='/activated/failed' element={<ActivatedFailedPage />} />
      </Routes>
    </main>
  );
}

export default App;
