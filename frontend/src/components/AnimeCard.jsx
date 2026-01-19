import React from 'react';
import { PlayCircle, Star } from 'lucide-react';
import { Link } from 'react-router-dom';

const AnimeCard = ({ title, thumb, endpoint, rating, episode, type }) => {
    // Determine if it's manga or anime based on endpoint structure or prop
    // Winbu endpoint usually just slug.
    // For now assume anime detail route.

    return (
        <Link to={`/detail/${endpoint.replace(/\/$/, "")}`} className="group relative block bg-secondary rounded-xl overflow-hidden shadow-lg border border-white/5 hover:border-primary/50 transition-all duration-300">
            {/* Image Container */}
            <div className="aspect-[2/3] overflow-hidden relative">
                <img
                    src={thumb}
                    alt={title}
                    className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                    loading="lazy"
                />

                {/* Overlay Gradient */}
                <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent opacity-60 group-hover:opacity-80 transition-opacity" />

                {/* Play Icon on Hover */}
                <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                    <PlayCircle className="w-12 h-12 text-primary fill-black/50" />
                </div>

                {/* Badges */}
                <div className="absolute top-2 left-2 flex flex-col gap-1">
                    {rating && (
                        <span className="px-2 py-1 bg-black/60 backdrop-blur-sm text-yellow-400 text-xs font-bold rounded flex items-center gap-1">
                            <Star className="w-3 h-3 fill-yellow-400" /> {rating}
                        </span>
                    )}
                </div>

                <div className="absolute top-2 right-2">
                    <span className="px-2 py-1 bg-primary/80 backdrop-blur-sm text-white text-xs font-bold rounded">
                        {type || 'Anime'}
                    </span>
                </div>
            </div>

            {/* Content */}
            <div className="p-3">
                <h3 className="text-white font-semibold text-sm line-clamp-2 leading-tight group-hover:text-primary transition-colors">
                    {title}
                </h3>
                {episode && (
                    <p className="text-gray-400 text-xs mt-1">{episode}</p>
                )}
            </div>
        </Link>
    );
};

export default AnimeCard;
