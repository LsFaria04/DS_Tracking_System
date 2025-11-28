import { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faLeaf, faInfoCircle, faTruck, faPlane, faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons';
import type { CarbonFootprintData } from '../utils/carbonFootprint';
import { getCarbonFootprintExplanation } from '../utils/carbonFootprint';

interface CarbonFootprintProps {
    data: CarbonFootprintData;
}

export default function CarbonFootprint({ data }: CarbonFootprintProps) {
    const [showTooltip, setShowTooltip] = useState(false);
    const [isOpen, setIsOpen] = useState(false);

    if (data.status === 'ERROR') {
        return (
            <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl">
                <button
                    id = "carb-button"
                    onClick={() => setIsOpen(!isOpen)}
                    className="w-full p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors rounded-2xl"
                >
                    <div className="flex items-center gap-3">
                        <FontAwesomeIcon icon={faLeaf} className="text-xl text-gray-400" />
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Carbon Footprint</h2>
                        <span className="text-xs text-gray-500 dark:text-gray-400">(Click to expand)</span>
                    </div>
                    <FontAwesomeIcon 
                        icon={isOpen ? faChevronUp : faChevronDown} 
                        className="text-gray-400"
                    />
                </button>
                {isOpen && (
                    <div className="px-6 pb-6">
                        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                            <p className="text-sm text-red-600 dark:text-red-400">
                                {data.message || 'Unable to calculate carbon footprint at this time.'}
                            </p>
                        </div>
                    </div>
                )}
            </div>
        );
    }

    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl">
            <button
                id = "carb-button"
                onClick={() => setIsOpen(!isOpen)}
                className="w-full p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors rounded-2xl"
            >
                <div className="flex items-center gap-3">
                    <FontAwesomeIcon icon={faLeaf} className="text-xl text-green-500" />
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Carbon Footprint</h2>
                    <div className="flex items-center gap-2">
                        <span className="text-2xl font-bold text-green-700 dark:text-green-400">
                            {data.totalEmissions.toFixed(2)}
                        </span>
                        <span className="text-sm text-gray-600 dark:text-gray-400">kg CO₂e</span>
                    </div>
                    <span className="text-xs text-gray-500 dark:text-gray-400">(Click to {isOpen ? 'collapse' : 'expand'})</span>
                </div>
                <FontAwesomeIcon 
                    icon={isOpen ? faChevronUp : faChevronDown} 
                    className="text-gray-400"
                />
            </button>

            {isOpen && (
                <div className="px-6 pb-6">
                    <div className="flex items-center justify-between mb-4">
                        <div></div>
                        <div className="relative">
                            <button
                                onMouseEnter={() => setShowTooltip(true)}
                                onMouseLeave={() => setShowTooltip(false)}
                                onClick={() => setShowTooltip(!showTooltip)}
                                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                            >
                                <FontAwesomeIcon icon={faInfoCircle} className="text-lg" />
                            </button>
                            {showTooltip && (
                                <div className="absolute right-0 top-8 w-80 bg-gray-900 dark:bg-gray-800 text-white text-xs rounded-lg p-3 shadow-lg z-10">
                                    <p className="font-semibold mb-2">How is this calculated?</p>
                                    <p>{getCarbonFootprintExplanation(data)}</p>
                                    <div className="absolute -top-1 right-4 w-2 h-2 bg-gray-900 dark:bg-gray-800 transform rotate-45"></div>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Main emission value */}
                    <div className="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 border border-green-200 dark:border-green-800 rounded-xl p-6 mb-4">
                        <div className="text-center">
                            <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">
                                Total CO₂ Emissions
                            </p>
                            <p className="text-5xl font-bold text-green-700 dark:text-green-400 mb-1">
                                {data.totalEmissions.toFixed(2)}
                            </p>
                            <p className="text-lg text-gray-600 dark:text-gray-400">
                                kg CO₂e
                            </p>
                            {data.status === 'ESTIMATED' && (
                                <p className="text-xs text-orange-600 dark:text-orange-400 mt-3 font-medium">
                                    Estimated Value
                                </p>
                            )}
                            {data.status === 'CALCULATED' && (
                                <p className="text-xs text-green-600 dark:text-green-400 mt-3 font-medium">
                                    Based on Actual Route
                                </p>
                            )}
                        </div>
                    </div>

                    {/* Breakdown */}
                    <div className="space-y-3">
                        {data.roadDistance > 0 && (
                            <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
                                <div className="flex items-center gap-3">
                                    <FontAwesomeIcon icon={faTruck} className="text-gray-600 dark:text-gray-400" />
                                    <div>
                                        <p className="text-sm font-medium text-gray-900 dark:text-white">Road Transport</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400">{data.roadDistance.toFixed(1)} km</p>
                                    </div>
                                </div>
                                <p className="text-sm font-semibold text-gray-900 dark:text-white">
                                    {data.roadEmissions.toFixed(2)} kg CO₂e
                                </p>
                            </div>
                        )}
                        
                        {data.airDistance > 0 && (
                            <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
                                <div className="flex items-center gap-3">
                                    <FontAwesomeIcon icon={faPlane} className="text-blue-600 dark:text-blue-400" />
                                    <div>
                                        <p className="text-sm font-medium text-gray-900 dark:text-white">Air Freight</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400">{data.airDistance.toFixed(1)} km</p>
                                    </div>
                                </div>
                                <p className="text-sm font-semibold text-gray-900 dark:text-white">
                                    {data.airEmissions.toFixed(2)} kg CO₂e
                                </p>
                            </div>
                        )}
                    </div>

                    {/* Context */}
                    <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-800">
                        <p className="text-xs text-gray-500 dark:text-gray-400 text-center">
                            {data.status === 'CALCULATED' 
                                ? 'Calculated based on the actual delivery route'
                                : 'Estimated based on route from seller to delivery location'}
                        </p>
                    </div>
                </div>
            )}
        </div>
    );
}
