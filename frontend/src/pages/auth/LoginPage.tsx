import React, { useState } from 'react';
import { Link } from 'react-router-dom';

const LoginPage: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Data Login (Dummy):', formData);
    alert('Tombol Sign In ditekan! Cek console.');
  };

  return (
    <div className="min-h-screen bg-gray-50 flex justify-center items-center p-4 sm:p-8">
      <div className="bg-white rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.08)] flex flex-col md:flex-row w-full max-w-[900px] border border-gray-200 overflow-hidden">
        <div className="hidden md:block w-1/2 p-3">
          <img
            src="https://images.unsplash.com/photo-1541701494587-cb58502866ab?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80"
            alt="Abstract Art"
            className="w-full h-full object-cover rounded-xl min-h-[500px]"
          />
        </div>

        <div className="w-full md:w-1/2 flex flex-col justify-center p-8 md:p-12">
          <div className="text-center mb-8">
            <h2 className="text-[28px] font-bold text-gray-800 mb-1">Welcome Back 👋</h2>
            <p className="text-sm text-gray-500">
              Tidak punya akun?{' '}
              <Link to="/register" className="text-[#3EA2E8] font-semibold hover:underline">
                Daftar
              </Link>
            </p>
          </div>

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
                className="w-full px-4 py-2.5 border border-gray-200 rounded-md text-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#1A4B69] focus:border-transparent transition-all"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-bold text-gray-700 mb-1.5" htmlFor="password">
                Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                placeholder="Password"
                value={formData.password}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-200 rounded-md text-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#1A4B69] focus:border-transparent transition-all"
                required
              />
            </div>

            <div className="flex justify-end pt-1 pb-3">
              <Link to="/forgot-password" className="text-sm text-[#3EA2E8] hover:underline">
                Forgot Password?
              </Link>
            </div>

            <div>
              <button
                type="submit"
                className="w-full bg-[#1A4B69] hover:bg-[#12364c] text-white text-sm font-semibold py-3 rounded-md transition-colors duration-200 shadow-sm"
              >
                Sign in
              </button>
            </div>
          </form>
        </div>
        
      </div>
    </div>
  );
};

export default LoginPage;