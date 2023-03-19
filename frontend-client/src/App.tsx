import { Routes, Route } from 'react-router-dom';
import UserActivationPage from './pages/UserActivationPage';

function App() {
  return (
    <main>
      <Routes>
        <Route path='/activate' element={<UserActivationPage />} />
      </Routes>
    </main>
  );
}

export default App;
