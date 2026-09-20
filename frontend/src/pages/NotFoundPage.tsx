import React from 'react';
import { useNavigate } from 'react-router-dom';

const NotFoundPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center items-center p-4">
      <div className="text-center max-w-md">
        <h1 className="text-9xl font-extrabold text-[#1A4B69] tracking-widest drop-shadow-sm">
          404
        </h1>
        <div className="bg-[#3EA2E8] text-white px-2 text-sm rounded rotate-12 absolute -mt-16 ml-32 shadow-sm">
          Page Not Found
        </div>
        
        <h2 className="mt-8 text-2xl font-bold text-gray-800">
          Waduh, nyasar!
        </h2>
        <p className="mt-4 text-gray-500 mb-8 leading-relaxed">
          Halaman yang Anda cari sepertinya udah dipindah, dihapus, atau emang gak pernah ada. Yuk balik ke jalan yang benar.
        </p>
        
        <button
          onClick={() => navigate('/')}
          className="px-8 py-3 bg-[#1A4B69] hover:bg-[#12364c] active:scale-95 text-white font-semibold rounded-xl transition-all shadow-sm flex items-center justify-center gap-2 mx-auto"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Kembali ke Beranda
        </button>
      </div>
    </div>
  );
};

export default NotFoundPage;