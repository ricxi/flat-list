import { useSearchParams } from 'react-router-dom';

function UserActivationPage() {
  const [searchParams] = useSearchParams();

  const handleClick = () => {
    const token = searchParams.get('token');
    if (token != null) {
      alert(token);
    }
  };

  return (
    <div style={{ textAlign: 'center' }}>
      <p>Please click on this button to activate your account.</p>
      <button onClick={handleClick}>activate</button>
    </div>
  );
}

export default UserActivationPage;
