import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Play } from 'lucide-react';
import { Link } from 'react-router-dom';

const HeroSlider = ({ slides }) => {
    const [current, setCurrent] = useState(0);

    useEffect(() => {
        const timer = setInterval(() => {
            setCurrent((prev) => (prev + 1) % slides.length);
        }, 5000);
        return () => clearInterval(timer);
    }, [slides]);

    if (!slides || slides.length === 0) return null;

    const nextSlide = () => setCurrent((c) => (c + 1) % slides.length);
    const prevSlide = () => setCurrent((c) => (c - 1 + slides.length) % slides.length);

    return (
        <div className="relative w-full h-[50vh] md:h-[60vh] overflow-hidden mb-8 group">
            {/* Background Image with Blur */}
            <div
                className="absolute inset-0 bg-cover bg-center blur-sm scale-110 opacity-50 transition-all duration-700"
                style={{ backgroundImage: `url(${slides[current].Thumb})` }}
            />
            <div className="absolute inset-0 bg-gradient-to-t from-background via-background/60 to-transparent" />
            <div className="absolute inset-0 bg-gradient-to-r from-background via-black/50 to-transparent" />

            {/* Content */}
            <div className="absolute inset-0 flex items-center">
                <div className="container mx-auto px-4 grid md:grid-cols-2 gap-8 items-center">
                    <div className="space-y-4 z-10 animate-fade-in-up">
                        <span className="text-primary font-bold tracking-widest text-sm uppercase">Featured</span>
                        <h1 className="text-4xl md:text-6xl font-bold text-white line-clamp-2 leading-tight">
                            {slides[current].Title}
                        </h1>
                        <div className="flex items-center gap-4 text-sm text-gray-300">
                            <span className="px-2 py-1 bg-yellow-500/20 text-yellow-400 border border-yellow-500/50 rounded">
                                {slides[current].Rating || 'N/A'}
                            </span>
                            <span>{slides[current].Status || 'Ongoing'}</span>
                            <span>{slides[current].Type || 'Anime'}</span>
                        </div>
                        <p className="text-gray-400 line-clamp-3 md:line-clamp-4 max-w-lg">
                            {/* Synopsis not always available in list, might need detail */}
                            Watch the latest episodes of {slides[current].Title} now on DramaBos Clone.
                        </p>
                        <div className="pt-4 flex gap-4">
                            <Link
                                to={`/detail/${slides[current].Endpoint.replace(/\/$/, "")}`}
                                className="flex items-center gap-2 px-8 py-3 bg-primary hover:bg-blue-600 text-white rounded-full font-bold transition-transform hover:scale-105 shadow-[0_0_20px_rgba(59,130,246,0.5)]"
                            >
                                <Play className="fill-white w-5 h-5" /> Watch Now
                            </Link>
                        </div>
                    </div>

                    {/* Poster Card (Hidden on mobile) */}
                    <div className="hidden md:block justify-self-end relative z-10">
                        <div className="w-64 aspect-[2/3] rounded-xl overflow-hidden shadow-[0_0_30px_rgba(0,0,0,0.5)] border border-white/10 transform rotate-3">
                            <img
                                src={slides[current].Thumb}
                                alt={slides[current].Title}
                                className="w-full h-full object-cover"
                            />
                        </div>
                    </div>
                </div>
            </div>

            {/* Arrows */}
            <button onClick={prevSlide} className="absolute left-4 top-1/2 -translate-y-1/2 p-2 bg-black/50 rounded-full hover:bg-primary text-white opacity-0 group-hover:opacity-100 transition-all">
                <ChevronLeft />
            </button>
            <button onClick={nextSlide} className="absolute right-4 top-1/2 -translate-y-1/2 p-2 bg-black/50 rounded-full hover:bg-primary text-white opacity-0 group-hover:opacity-100 transition-all">
                <ChevronRight />
            </button>

            {/* Dots */}
            <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                {slides.map((_, idx) => (
                    <button
                        key={idx}
                        onClick={() => setCurrent(idx)}
                        className={`w-2 h-2 rounded-full transition-all ${current === idx ? 'w-8 bg-primary' : 'bg-gray-500'}`}
                    />
                ))}
            </div>
        </div>
    );
};

export default HeroSlider;
