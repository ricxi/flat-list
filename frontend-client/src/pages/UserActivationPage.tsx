import { useNavigate, useSearchParams } from 'react-router-dom';

const USER_ACTIVATION_ENDPOINT = 'http://localhost:9000/v1/user/activate';

function UserActivationPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const handleClick = async () => {
    const token = searchParams.get('token');

    if (token === null) {
      navigate('/activated/failed');
      return;
    }

    // Should I put the token in the body?
    const res = await fetch(USER_ACTIVATION_ENDPOINT + '/' + token, {
      method: 'PUT',
    });

    if (res.status !== 200) {
      navigate('/activated/failed');
      return;
    }

    navigate('/activated/success');
  };

  return (
    <div style={{ textAlign: 'center' }}>
      <p>Please click on this button to activate your account.</p>
      <button onClick={handleClick}>activate</button>
    </div>
  );
}

export default UserActivationPage;
