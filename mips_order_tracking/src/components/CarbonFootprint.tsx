import { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faLeaf, faInfoCircle, faTruck, faPlane, faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons';
import type { CarbonFootprintData } from '../utils/carbonFootprint';
import { getCarbonFootprintExplanation } from '../utils/carbonFootprint';
import {
  Box,
  Button,
  Typography,
  Divider,
  Tooltip,
  useTheme,
} from '@mui/material';

interface CarbonFootprintProps {
    data: CarbonFootprintData;
}

export default function CarbonFootprint({ data }: CarbonFootprintProps) {
    const [showTooltip, setShowTooltip] = useState(false);
    const [isOpen, setIsOpen] = useState(false);
    const theme = useTheme();

    if (data.status === 'ERROR') {
        return (
            <Box sx={{ border: `1px solid ${theme.palette.divider}`, borderRadius: 1 }}>
                <Button
                    id="carb-button"
                    onClick={() => setIsOpen(!isOpen)}
                    sx={{
                        width: '100%',
                        p: 2,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        textTransform: 'none',
                        color: 'text.primary',
                        '&:hover': { bgcolor: 'action.hover' }
                    }}
                >
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flex: 1 }}>
                        <FontAwesomeIcon icon={faLeaf} style={{ color: theme.palette.error.main, fontSize: '1.25rem' }} />
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>Carbon Footprint</Typography>
                        <Typography variant="caption" color="textSecondary">(Click to expand)</Typography>
                    </Box>
                    <FontAwesomeIcon icon={isOpen ? faChevronUp : faChevronDown} />
                </Button>
                {isOpen && (
                    <Box sx={{ p: 2, bgcolor: 'rgba(211, 47, 47, 0.05)', borderTop: `1px solid ${theme.palette.divider}` }}>
                        <Typography sx={{ color: theme.palette.error.main }}>
                            {data.message || 'Unable to calculate carbon footprint at this time.'}
                        </Typography>
                    </Box>
                )}
            </Box>
        );
    }

    return (
        <Box sx={{ border: `1px solid ${theme.palette.divider}`, borderRadius: 1 }}>
            <Button
                id="carb-button"
                onClick={() => setIsOpen(!isOpen)}
                sx={{
                    width: '100%',
                    p: 2,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    textTransform: 'none',
                    color: 'text.primary',
                    '&:hover': { bgcolor: 'action.hover' }
                }}
            >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flex: 1 }}>
                    <FontAwesomeIcon icon={faLeaf} style={{ color: theme.palette.success.main, fontSize: '1.25rem' }} />
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>Carbon Footprint</Typography>
                    {!isOpen && (
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                            <Typography variant="h5" sx={{ fontWeight: 700, color: theme.palette.success.main }}>
                                {data.totalEmissions.toFixed(2)}
                            </Typography>
                            <Typography variant="caption" color="textSecondary">kg CO₂e</Typography>
                        </Box>
                    )}
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="caption" color="textSecondary">
                        {isOpen ? '(Click to collapse)' : '(Click to expand)'}
                    </Typography>
                    <FontAwesomeIcon icon={isOpen ? faChevronUp : faChevronDown} />
                </Box>
            </Button>

            {isOpen && (
                <Box sx={{ p: 2, borderTop: `1px solid ${theme.palette.divider}` }}>
                    {/* Info Button */}
                    <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 2 }}>
                        <Tooltip
                            open={showTooltip}
                            onOpen={() => setShowTooltip(true)}
                            onClose={() => setShowTooltip(false)}
                            title={
                                <Box sx={{ whiteSpace: 'normal', maxWidth: 300 }}>
                                    <Typography variant="body2" sx={{ fontWeight: 600, mb: 1 }}>How is this calculated?</Typography>
                                    <Typography variant="caption">{getCarbonFootprintExplanation(data)}</Typography>
                                </Box>
                            }
                        >
                            <Button
                                size="small"
                                onClick={() => setShowTooltip(!showTooltip)}
                                sx={{ minWidth: 'auto', p: 0.5, color: 'text.secondary' }}
                            >
                                <FontAwesomeIcon icon={faInfoCircle} />
                            </Button>
                        </Tooltip>
                    </Box>

                    {/* Total Emissions Box */}
                    <Box
                        sx={{
                            textAlign: 'center',
                            p: 3,
                            mb: 2,
                            bgcolor: 'rgba(76, 175, 80, 0.05)',
                            border: `1px solid ${theme.palette.success.light}`,
                            borderRadius: 1,
                        }}
                    >
                        <Typography variant="body2" sx={{ color: 'text.secondary', mb: 1 }}>
                            Total CO₂ Emissions
                        </Typography>
                        <Typography variant="h3" sx={{ fontWeight: 700, color: theme.palette.success.main, mb: 0.5 }}>
                            {data.totalEmissions.toFixed(2)}
                        </Typography>
                        <Typography variant="body2" sx={{ color: 'text.secondary', mb: 2 }}>
                            kg CO₂e
                        </Typography>
                        {data.status === 'ESTIMATED' && (
                            <Typography variant="caption" sx={{ color: theme.palette.warning.main, fontWeight: 600 }}>
                                Estimated Value
                            </Typography>
                        )}
                        {data.status === 'CALCULATED' && (
                            <Typography variant="caption" sx={{ color: theme.palette.success.main, fontWeight: 600 }}>
                                Based on Actual Route
                            </Typography>
                        )}
                    </Box>

                    {/* Breakdown Items */}
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                        {data.roadDistance > 0 && (
                            <Box
                                sx={{
                                    display: 'flex',
                                    justifyContent: 'space-between',
                                    alignItems: 'center',
                                    p: 1.5,
                                    bgcolor: 'action.hover',
                                    borderRadius: 1,
                                }}
                            >
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <FontAwesomeIcon icon={faTruck} />
                                    <Box>
                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>Road Transport</Typography>
                                        <Typography variant="caption" color="textSecondary">{data.roadDistance.toFixed(1)} km</Typography>
                                    </Box>
                                </Box>
                                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                    {data.roadEmissions.toFixed(2)} kg CO₂e
                                </Typography>
                            </Box>
                        )}

                        {data.airDistance > 0 && (
                            <Box
                                sx={{
                                    display: 'flex',
                                    justifyContent: 'space-between',
                                    alignItems: 'center',
                                    p: 1.5,
                                    bgcolor: 'action.hover',
                                    borderRadius: 1,
                                }}
                            >
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <FontAwesomeIcon icon={faPlane} />
                                    <Box>
                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>Air Freight</Typography>
                                        <Typography variant="caption" color="textSecondary">{data.airDistance.toFixed(1)} km</Typography>
                                    </Box>
                                </Box>
                                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                    {data.airEmissions.toFixed(2)} kg CO₂e
                                </Typography>
                            </Box>
                        )}
                    </Box>

                    <Divider sx={{ my: 2 }} />
                    <Typography variant="caption" color="textSecondary" sx={{ display: 'block', textAlign: 'center' }}>
                        {data.status === 'CALCULATED' 
                            ? 'Calculated based on the actual delivery route'
                            : 'Estimated based on route from seller to delivery location'}
                    </Typography>
                </Box>
            )}
        </Box>
    );
}
