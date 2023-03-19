import { useSearchParams } from 'react-router-dom';

function UserActivationPage() {
  const [searchParams] = useSearchParams();

  const handleClick = async () => {
    const token = searchParams.get('token');

    if (token === null) {
      alert('missing token');
    }

    const url = `http://localhost:9000/v1/user/activate/${token}`;

    const res = await fetch(url, {
      method: 'PUT',
    });

    console.log(token);

    if (res.status !== 200) {
      console.error('problem activatign your account, please try again');
    } else {
      console.log('success!');
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
