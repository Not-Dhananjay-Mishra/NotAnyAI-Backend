package router

import codingmodel "server/models/CodingModel"

func SusTemp() codingmodel.PostCodeResponse {
	// Initialize maps properly
	sus := make(map[string]string)
	eee := make(map[string]string)

	sus["App.js"] = `
import React from 'react';
import HomePage from './index.js';

function App() {
  return (
    <div className="App">
      <HomePage />
    </div>
  );
}

export default App;
`

	sus["index.js"] = `
import React from 'react';
import { motion } from 'framer-motion';

const HomePage = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-800 via-fuchsia-700 to-pink-600 flex flex-col items-center justify-center text-center p-4">
      <motion.div
        initial={{ opacity: 0, y: -50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8 }}
        className="max-w-4xl"
      >
        <motion.h1
          className="text-6xl md:text-7xl font-extrabold text-white tracking-wide drop-shadow-lg leading-tight"
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.2 }}
        >
          Unveil the Unseen. Discover the Unexpected.
        </motion.h1>

        <motion.p
          className="mt-6 text-lg md:text-xl text-purple-100 max-w-2xl mx-auto opacity-90"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
        >
          Dive into a world where curiosity leads the way. Our app reveals insights you never knew existed.
        </motion.p>

        <motion.button
          className="mt-10 px-10 py-4 bg-fuchsia-500 text-white rounded-full text-xl font-semibold hover:bg-fuchsia-400 transition-colors duration-300 focus:outline-none focus:ring-4 focus:ring-fuchsia-300 focus:ring-opacity-70 shadow-lg"
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.6, delay: 0.6 }}
          aria-label="Get Started with our app"
        >
          Get Started Now
        </motion.button>
      </motion.div>
    </div>
  );
};

export default HomePage;
`

	return codingmodel.PostCodeResponse{FrontendCode: sus, BackendCode: eee}
}
