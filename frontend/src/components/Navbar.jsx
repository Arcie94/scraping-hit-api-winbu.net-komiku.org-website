import React from 'react';
import { Code2, Github } from 'lucide-react';
import { Link } from 'react-router-dom';

const Navbar = () => {
    return (
        <nav className="sticky top-0 z-50 bg-background/80 backdrop-blur-md border-b border-white/10">
            <div className="container mx-auto px-4">
                <div className="flex items-center justify-between h-16">
                    {/* Logo */}
                    <Link to="/" className="flex items-center space-x-2">
                        <Code2 className="w-8 h-8 text-primary" />
                        <span className="text-xl font-bold bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
                            API Explorer
                        </span>
                    </Link>

                    {/* Right Menu */}
                    <div className="flex items-center space-x-6">
                        {/* GitHub Link */}
                        <a
                            href="https://github.com/Arcie94"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg transition-all duration-300 group"
                        >
                            <Github className="w-5 h-5 text-gray-400 group-hover:text-white transition-colors" />
                            <span className="text-gray-300 group-hover:text-white font-medium">GitHub</span>
                        </a>
                        <span className="px-3 py-1 bg-green-500/20 text-green-400 border border-green-500/50 rounded-full text-xs font-bold">
                            ONLINE
                        </span>
                    </div>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
