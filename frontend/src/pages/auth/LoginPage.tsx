import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../store/hooks'; 
import { login } from '../../features/auth/authSlice';
import toast from 'react-hot-toast';

const useLoginController = () => {
  const navigate = useNavigate();

  const { isAuthenticated } = useAppSelector((state) => state.auth);
  useEffect(() => {
    if (isAuthenticated) {
      navigate("/", { replace: true });
    }
  }, [isAuthenticated, navigate]);
  const dispatch = useAppDispatch();

  const { loading, error } = useAppSelector((state) => state.auth);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await dispatch(login(formData)).unwrap();
      toast.success('Login berhasil!');
      navigate('/');
    } catch (err) {
      toast.error(err as string || 'Login gagal');
    }
  };

  return {
    formData,
    loading,
    error,
    handleChange,
    handleSubmit
  };
};

const LoginPage: React.FC = () => {
  const { formData, loading, handleChange, handleSubmit } = useLoginController();
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="min-h-screen bg-gray-50 flex justify-center items-center p-4 sm:p-8">
      <div className="bg-white rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.08)] flex flex-col md:flex-row w-full max-w-[900px] border border-gray-200 overflow-hidden">
        <div className="hidden md:block w-1/2 p-3">
          <img
            src="./src/assets/bid_img.webp"
            alt="Auction Picture"
            className="w-full h-full object-cover rounded-xl min-h-[500px]"
          />
        </div>

        <div className="w-full md:w-1/2 flex flex-col justify-center p-8 md:p-12">
          <div className="text-center mb-8">
            <h2 className="text-[28px] font-bold text-gray-800 mb-1">Welcome Back 👋</h2>
            {/* <p className="text-sm text-gray-500">
              Tidak punya akun?{' '}
              <Link to="/register" className="text-[#3EA2E8] font-semibold hover:underline">
                Daftar
              </Link>
            </p> */}
          </div>

          {/* {error && (
            <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 text-sm rounded">
              {error}
            </div>
          )} */}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-1.5" htmlFor="username">
                Username
              </label>
              <input
                type="text"
                id="username"
                name="username"
                placeholder="Username"
                value={formData.username}
                onChange={handleChange}
                autoFocus
                className="w-full px-4 py-2.5 border border-gray-200 rounded-md text-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#1A4B69] focus:border-transparent transition-all"
                required
              />
            </div>

            <div>
                <label className="block text-sm font-bold text-gray-700 mb-1.5" htmlFor="password">
                  Password
                </label>
                
                <div className="relative flex items-center">
                  
                  <input
                    type={showPassword ? "text" : "password"}
                    id="password"
                    name="password"
                    placeholder="Password"
                    value={formData.password}
                    onChange={handleChange}
                    className="w-full px-4 py-2.5 pr-12 border border-gray-200 rounded-md text-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#1A4B69] focus:border-transparent transition-all"
                    required
                  />

                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 p-1.5 text-gray-400 hover:text-gray-600 focus:outline-none rounded-full hover:bg-gray-100 transition-colors"
                    aria-label={showPassword ? "Hide password" : "Show password"}
                  >
                    {showPassword ? (
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth="2">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path strokeLinecap="round" strokeLinejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                    ) : (
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth="2">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                      </svg>
                    )}
                  </button>
                  
                </div>
            </div>

            {/* <div className="flex justify-end pt-1 pb-3">
              <Link to="/forgot-password" className="text-sm text-[#3EA2E8] hover:underline">
                Forgot Password?
              </Link>
            </div> */}

            <div>
              <button
                type="submit"
                disabled={loading}
                className={`w-full text-white text-sm font-semibold py-3 rounded-md transition-colors duration-200 shadow-sm ${
                  loading ? 'bg-gray-400 cursor-not-allowed' : 'bg-[#1A4B69] hover:bg-[#12364c]'
                }`}
              >
                {loading ? 'Signing in...' : 'Sign in'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;