import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const OAuthCallback: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const url = new URL(window.location.href);
    const token = url.searchParams.get('token');
    const email = url.searchParams.get('email');

    console.log('OAuth callback params:', { token, email });

    if (token && email) {
      localStorage.setItem('access_token', token);
      localStorage.setItem(
        'user',
        JSON.stringify({
          id: '',
          username: email.split('@')[0],
          email,
          role: 'user',
          is_verified: true,
          created_at: '',
          updated_at: '',
        })
      );

      window.history.replaceState({}, document.title, '/dashboard');

      navigate('/dashboard', { replace: true });
    } else {
      console.warn('OAuth callback missing token/email, redirecting to login');
      navigate('/login', { replace: true });
    }
  }, [navigate]);

  return "Logging you in...";
};

export default OAuthCallback;